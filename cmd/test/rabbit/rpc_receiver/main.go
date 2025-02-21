package main

import (
	"context"
	"onx-outgoing-go/internal/pkg/logger"
	"onx-outgoing-go/internal/pkg/rabbitmq"
	"onx-outgoing-go/internal/service/worker"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger.Setup()
	rb, err := rabbitmq.NewConnectionManager(ctx, &rabbitmq.Config{
		Username: "test",
		Password: "test",
		Host:     "localhost",
		Port:     5672,
	})
	if err != nil {
		panic(err)
	}

	publisher, err := rabbitmq.NewPublisher(ctx, rb)
	if err != nil {
		panic(err)
	}

	s, err := worker.NewService(ctx, rb, publisher)
	if err != nil {
		panic(err)
	}

	if err := s.SampleMessageReceiverRPCRabbit(); err != nil {
		logger.Error.Println(err)
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()
}
