package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"release-platform/internal/api"
	"release-platform/internal/config"
	"release-platform/internal/db"
	"release-platform/internal/service"

	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.StorageRoot, 0o755); err != nil {
		log.Fatalf("create storage root: %v", err)
	}

	database, err := db.New(cfg.DatabaseDriver, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("init database: %v", err)
	}

	var queueClient *asynq.Client
	if strings.EqualFold(cfg.ScanMode, "async") {
		redisOpt := asynq.RedisClientOpt{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		}
		queueClient = asynq.NewClient(redisOpt)
		defer queueClient.Close()
	} else {
		log.Printf("scan mode is sync, Redis queue client is disabled")
	}

	artifactSvc := service.NewArtifactService(cfg.StorageRoot)
	scannerSvc := service.NewScannerService(database)
	releaseSvc := service.NewReleaseService(database, artifactSvc, scannerSvc, queueClient, cfg.QueueName)

	handler := api.NewHandler(releaseSvc, artifactSvc)
	router := api.NewRouter(handler, "web")

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("HTTP server started on %s", cfg.HTTPAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve http: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down HTTP server")
	if err := httpServer.Close(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	fmt.Println("server stopped")
}
