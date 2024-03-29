package storage

import (
	"sync"

	"github.com/kyle8615/ethereum-parser/v1/internal/model"
)

type Storage interface {
	PutTransaction(address string, txn model.Transaction) error
	GetTransactionByAddress(address string) ([]model.Transaction, error)
	SubscribeAddress(address string) error
	IsAddressSubscribed(address string) bool
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		mu:                &sync.RWMutex{},
		txMap:             make(map[string][]model.Transaction),
		subscribedAddress: make(map[string]struct{}),
	}
}

type memoryStorage struct {
	mu                *sync.RWMutex
	txMap             map[string][]model.Transaction
	subscribedAddress map[string]struct{}
}

func (m *memoryStorage) PutTransaction(address string, txn model.Transaction) error {
	m.mu.Lock()
	if _, ok := m.txMap[address]; !ok {
		m.txMap[address] = []model.Transaction{}
	}

	m.txMap[address] = append(m.txMap[address], txn)
	m.mu.Unlock()

	return nil
}

func (m *memoryStorage) GetTransactionByAddress(address string) ([]model.Transaction, error) {
	m.mu.RLock()
	txns := m.txMap[address]
	m.mu.RUnlock()

	return txns, nil
}

func (m *memoryStorage) SubscribeAddress(address string) error {
	m.mu.RLock()
	_, ok := m.subscribedAddress[address]
	m.mu.RUnlock()

	if !ok {
		m.mu.Lock()
		m.subscribedAddress[address] = struct{}{}
		m.mu.Unlock()
	}

	return nil
}

func (m *memoryStorage) IsAddressSubscribed(address string) bool {
	m.mu.RLock()
	_, exist := m.subscribedAddress[address]
	m.mu.RUnlock()
	return exist
}
