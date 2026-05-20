package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

func main() {
	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379"}
	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
	})
	mux := asynq.NewServeMux()
	mux.HandleFunc("printer:receipt", handlePrintReceipt)
	mux.HandleFunc("printer:kitchen", handlePrintKitchen)
	if err := server.Run(mux); err != nil {
		log.Fatalf("Could not run worker server: %v", err)
	}
}

func handlePrintReceipt(ctx context.Context, t *asynq.Task) error {
	fmt.Printf("[WORKER] Printing receipt: %s\n", string(t.Payload()))
	return nil
}

func handlePrintKitchen(ctx context.Context, t *asynq.Task) error {
	fmt.Printf("[WORKER] Printing kitchen ticket: %s\n", string(t.Payload()))
	return nil
}