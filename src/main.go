package main

import (
	"app/src/app_settings"
	"app/src/core/database"
	"app/src/services/brokers"
	"app/src/services/grpc"
	"app/src/services/worker"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const consumerRestartDelay = 5 * time.Second

func main() {
	configs, err := app_settings.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.ConnectGORM(ctx, configs.Postgres.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}()

	broker := brokers.NewKafkaBroker(configs.Brokers.KafkaSettings)
	defer func() {
		if err := broker.Close(); err != nil {
			log.Printf("Failed to close broker: %v", err)
		}
	}()

	dependencies := worker.NewDependencies(db)

	runner := brokers.NewRunner(broker, consumerRestartDelay)
	worker.InitConsumers(runner, dependencies)

	grpcServer, err := grpc.NewServer(configs.Grpc)
	if err != nil {
		log.Fatalf("Failed to create gRPC server: %v", err)
	}
	worker.InitGrpc(grpcServer, dependencies)

	runnerDone := make(chan struct{})
	go func() {
		runner.Run(ctx)
		close(runnerDone)
	}()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC server starting on %s", grpcServer.Address())
		if err := grpcServer.Serve(); err != nil {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Print("Shutdown signal received")
	case err := <-serverErr:
		log.Printf("gRPC server stopped: %v", err)
	}

	stop()
	grpcServer.GracefulStop()
	<-runnerDone
}
