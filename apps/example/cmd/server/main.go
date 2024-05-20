package main

import (
	"context"
	"k8s-devenv/example/internal/config"
	"k8s-devenv/example/internal/handlers"
	"os"
	"os/signal"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	go cfg.Kafka.Consume(ctx, handlers.GetMessageHandler(cfg))
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
	cfg.Logger.Info("Shutting down")
	done := make(chan struct{})
	go func() {
		defer close(done)
		cfg.Kafka.Client.Close()
	}()

	select {
	case <-sigs:
		cfg.Logger.Info("Force shutdown")
	case <-done:
	}
}
