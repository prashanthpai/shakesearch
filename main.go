package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/prashanthpai/shakesearch/internal/api"
	"github.com/prashanthpai/shakesearch/pkg/server"
	"github.com/prashanthpai/shakesearch/pkg/shake"
)

const (
	inputFile = "completeworks.txt"
)

func main() {
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("os.Open(%s) failed: %s", inputFile, err.Error())
	}
	defer f.Close()

	// create instance of searcher
	shakeSearcher := shake.NewSearcher()
	if err := shakeSearcher.Load(f); err != nil {
		log.Fatalf("shakeSearcher.Load() failed: %s", err.Error())
	}

	// create HTTP handlers
	handlers := api.NewHandler(shakeSearcher)

	// configure server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	srv := server.New(&server.Config{
		Addr: fmt.Sprintf(":%s", port),
	}, handlers)

	// start server
	if err := srv.Start(); err != nil {
		log.Fatalf("server.Start() failed: %s", err.Error())
	}

	// handle shutdown gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	log.Println("Received interrupt. Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("http.Server.Shutdown() failed: %s\n", err)
	}
}
