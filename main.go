package main

import (
	"context"
	"go-masstransit-interop/internal/configuration"
	"go-masstransit-interop/internal/rabbitmq"
	"go-masstransit-interop/internal/worker"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg, err := configuration.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	workerPool := worker.NewWorkerPool(cfg.WorkerPool.Size)
	workerPool.Start()
	defer workerPool.Stop()

	registerTypes()
	registerHandlers()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQ.URL, cfg.RabbitMQ.QueueName, workerPool)
	if err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer: %s", err)
	}
	defer consumer.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := consumer.StartConsuming(ctx); err != nil {
			log.Printf("Error in consuming messages: %s", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received")

	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	workerPool.Stop()
	wg.Wait()
	log.Println("Application gracefully stopped")
}
