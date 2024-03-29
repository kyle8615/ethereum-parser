package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kyle8615/ethereum-parser/v1/internal/ethereum"
	"github.com/kyle8615/ethereum-parser/v1/internal/parser"
	"github.com/kyle8615/ethereum-parser/v1/internal/storage"
	"github.com/kyle8615/ethereum-parser/v1/pkg/api"
)

func main() {
	ctxApp, cancelApp := context.WithCancel(context.Background())

	client := ethereum.NewClient("https://cloudflare-eth.com")
	store := storage.NewMemoryStorage()
	p, err := parser.NewParser(ctxApp, client, store)
	if err != nil {
		panic("initiate parser failed")
	}

	server := api.NewServer(p)

	ch := make(chan struct{})
	go func() {
		defer func() {
			ch <- struct{}{}
		}()

		fmt.Println("Service Start")
		if err := server.Start(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("Service End")
	}()

	// Setting up signal catching
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	cancelApp()
	fmt.Println("Shutting down server...")

	// set waiting time for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Failed to shutdown server: ", err)
	}

	<-ch // Wait for the server goroutine to complete
	fmt.Println("Server gracefully stopped")
}
