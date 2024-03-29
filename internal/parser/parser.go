package parser

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kyle8615/ethereum-parser/v1/internal/ethereum"
	"github.com/kyle8615/ethereum-parser/v1/internal/model"
	"github.com/kyle8615/ethereum-parser/v1/internal/storage"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []model.Transaction
}

type parser struct {
	client       *ethereum.Client
	storage      storage.Storage
	parsedHeight int
	mu           *sync.RWMutex
}

func NewParser(ctx context.Context, client *ethereum.Client, storage storage.Storage) (Parser, error) {

	if client == nil {
		return nil, fmt.Errorf("ehtereum client is essential for parser initialization")
	}

	if storage == nil {
		return nil, fmt.Errorf("storage is essential for parser initialization")
	}

	instance := &parser{
		client:       client,
		storage:      storage,
		mu:           &sync.RWMutex{},
		parsedHeight: 0x0,
	}

	// To follow the asked condition so invoke init here.
	// but the init function not only set configuration also do network request
	// and change status, it might be better to call it separately from the constructor
	// to give users more control over when those side effects occur.
	err := instance.init(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init parser")
	}

	return instance, nil
}

func (p *parser) GetCurrentBlock() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.parsedHeight
}

func (p *parser) Subscribe(address string) bool {
	// XXX: check the address is valid for the Ethereum
	err := p.storage.SubscribeAddress(strings.ToLower(address))
	return err == nil
}

func (p *parser) GetTransactions(address string) []model.Transaction {
	// omit error
	txns, _ := p.storage.GetTransactionByAddress(address)
	return txns
}

func (p *parser) init(ctx context.Context) error {

	// load processed height from persistent storage. Reset for the accessment

	go func() {
		for {

			select {
			case <-ctx.Done():
				return
			default:

				blockHeight, err := p.client.GetLatestBlockNumber()
				if err != nil {
					fmt.Println("get latest block number failed", err)
					// XXX: naively waiting for 3s to retry for now. use consisent increasing retry and monitor event.
					time.Sleep(3 * time.Second)
					continue
				}

				if blockHeight <= p.parsedHeight {
					fmt.Println("Caught up with the latest block, slow down.")
					time.Sleep(20 * time.Second)
					continue
				}

				block, err := p.client.GetBlockByNumber(p.parsedHeight + 1)
				if err != nil {
					fmt.Println("get block by number error", err)
					// XXX: naively waiting for 3s to retry for now. use consisent increasing retry and monitor event.
					time.Sleep(3 * time.Second)
					continue
				}

				// iterate block transactions to store subscribed incoming and outgoing tx
				for _, tx := range block.Transactions {
					if p.storage.IsAddressSubscribed(strings.ToLower(tx.From)) {
						// omit put error, if error occur in production env put error to a queue and try later
						p.storage.PutTransaction(strings.ToLower(tx.From), tx)
					}

					if p.storage.IsAddressSubscribed(strings.ToLower(tx.To)) {
						// omit put error, if error occur in production env put error to a queue and try later
						p.storage.PutTransaction(strings.ToLower(tx.To), tx)
					}

				}

				p.mu.Lock()
				p.parsedHeight++
				p.mu.Unlock()

				fmt.Println("proceed height: ", p.parsedHeight)
				// Update the processed height in persistent storage in real env.
			}
		}
	}()

	return nil
}
