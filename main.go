package main

import (
	"context"
	"fmt"
	"net/http"
	"onx-outgoing-go/config"
	"onx-outgoing-go/internal"
	"onx-outgoing-go/internal/pkg/logger"
	"onx-outgoing-go/internal/pkg/rabbitmq"
	"onx-outgoing-go/internal/pkg/redis"
	"onx-outgoing-go/internal/pkg/temporal"
	"onx-outgoing-go/internal/pkg/validation"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// @title ONX Livechat API
// @version 2.0
// @description API documentation for the ONX Livechat service.
// @BasePath /api
// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @type apiKey
func main() {
	env, err := config.GetEnv()
	if err != nil {
		panic(err)
	}
	logger.Setup()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	rds, err := setupRedis(ctx, env)
	if err != nil {
		logger.Error.Println("Error connecting to redis")
		panic(err)
	}

	rb, err := setupRabbitMQ(ctx, env)
	if err != nil {
		logger.Error.Println("Error connecting to redis")
		panic(err)
	}

	publisher, err := rabbitmq.NewPublisher(ctx, rb)
	if err != nil {
		panic(err)
	}

	// Initialize Temporal connection early and defer closing it until the app shuts down.
	err = temporal.InitTemporalConnection(ctx)
	if err != nil {
		fmt.Println("Unable to create Temporal client", err)
		panic(err)
	}
	// defer temporal.CloseTemporalConnection()

	setupServer(&config.SetupServerDto{
		Ctx:    &ctx,
		Cancel: cancel,
		Wg:     &wg,
		Env:    env,
		Rds:    rds,
		Rb:     rb,
		Pb:     publisher,
		Tp:     temporal.GetTemporalClient(),
	})

	// At this point, the server shutdown sequence has completed.
	// Give workers a moment to cleanup (if needed) then close the Temporal connection.
	time.Sleep(2 * time.Second) // optional: adjust for proper graceful shutdown
	temporal.CloseTemporalConnection()
}

func setupRedis(ctx context.Context, env *config.Config) (redis.IRedis, error) {
	return redis.Setup(ctx, &redis.Config{
		Host:     env.RedisHost,
		Port:     env.RedisPort,
		Password: env.RedisPass,
		PoolSize: env.RedisPoolSize,
	})
}

func setupRabbitMQ(ctx context.Context, env *config.Config) (*rabbitmq.ConnectionManager, error) {
	rb, err := rabbitmq.NewConnectionManager(ctx, &rabbitmq.Config{
		Username: env.RabbitUser,
		Password: env.RabbitPass,
		Host:     env.RabbitHost,
		Port:     env.RabbitPort,
	})
	if err != nil {
		panic(err)
	}

	rb.InitRPCClient()

	return rb, err
}

func setupServer(payload *config.SetupServerDto) {
	rds := payload.Rds
	env := payload.Env
	ctx := payload.Ctx
	cancel := payload.Cancel
	wg := payload.Wg
	rb := payload.Rb
	pb := payload.Pb
	tp := payload.Tp

	defer func() {
		if rds != nil {
			_ = rds.Close()
		}
		cancel()
		wg.Wait()
	}()

	err := validation.Setup()
	if err != nil {
		logger.Error.Println("Failed to setup validation")
		panic(err)
	}

	e := gin.Default()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", env.AppPort),
		Handler: e,
		//IdleTimeout:  1 * time.Minute,
		//ReadTimeout:  1 * time.Minute,
		//WriteTimeout: 1 * time.Minute,
	}

	internal.Setup(*env, e, *ctx, wg, rds, rb, pb, tp)

	go func() {
		logger.HTTP.Println("========= Server Started =========")
		logger.HTTP.Println("=========", env.AppPort, "=========")
		if err := server.ListenAndServe(); err != nil {
			logger.HTTP.Println(err)
			logger.HTTP.Println("========= Server Ended =========")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done
	logger.Error.Println("Received terminate, graceful shutdown", sigChan)
	_ = server.Shutdown(*ctx)
}
