package main

import (
	"log"
	"strings"

	"release-platform/internal/config"
	"release-platform/internal/db"
	"release-platform/internal/queue"
	"release-platform/internal/service"

	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.Load()

	if strings.EqualFold(cfg.ScanMode, "sync") {
		log.Println("scan mode is sync, worker is not required")
		return
	}

	database, err := db.New(cfg.DatabaseDriver, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("init database: %v", err)
	}

	artifactSvc := service.NewArtifactService(cfg.StorageRoot)
	scannerSvc := service.NewScannerService(database)
	releaseSvc := service.NewReleaseService(database, artifactSvc, scannerSvc, nil, cfg.QueueName)

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		},
		asynq.Config{
			Concurrency: 6,
			Queues: map[string]int{
				cfg.QueueName: 10,
			},
		},
	)

	mux := queue.NewMux(releaseSvc)
	log.Printf("worker started, listening queue=%s", cfg.QueueName)
	if err := server.Run(mux); err != nil {
		log.Fatalf("run worker: %v", err)
	}
}
