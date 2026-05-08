package main

import (
	"log"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.LoadConfig()

	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       0,
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	// Register task handlers here
	// mux.HandleFunc("sync:journal", worker.HandleSyncJournalTask)

	log.Println("Starting background worker...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not start worker: %v", err)
	}
}
