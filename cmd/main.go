package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"service-courier/internal/event"
	"service-courier/internal/factory"
	"service-courier/internal/gateway/order"
	"service-courier/internal/handler"
	"service-courier/internal/middleware"
	"service-courier/internal/repository"
	"service-courier/internal/service"
	"service-courier/internal/transport/kafka"
	"service-courier/pkg/config"
	"service-courier/pkg/database"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandler(cancel)

	dbPool, err := database.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	courierRepo := repository.NewCourierRepository(dbPool)
	deliveryRepo := repository.NewDeliveryRepository(dbPool)

	timeFactory := factory.NewDeliveryTimeFactory()

	courierService := service.NewCourierService(courierRepo)
	deliveryService := service.NewDeliveryService(deliveryRepo, courierRepo, timeFactory)

	backgroundService := service.NewBackgroundService(deliveryRepo, 10*time.Second)
	go backgroundService.StartOverdueCheck(ctx)

	orderGateway, err := order.NewGateway(cfg.OrderServiceURL)
	if err != nil {
		log.Fatalf("Failed to create order gateway: %v", err)
	}
	defer func() {
		if err := orderGateway.Close(); err != nil {
			slog.Error("Failed to close order gateway", "error", err)
		}
	}()

	orderPoller := service.NewOrderPoller(orderGateway, deliveryService)
	go orderPoller.Start(ctx)
	slog.Info("Order poller started", "interval", "5 seconds")

	eventHandlerFactory := event.NewEventHandlerFactory(deliveryService, orderGateway)
	eventProcessor := event.NewEventProcessor(eventHandlerFactory)

	brokers := []string{cfg.Kafka.Broker}
	kafkaConsumer, err := kafka.NewConsumer(
		brokers,
		cfg.Kafka.Topic,
		cfg.Kafka.ConsumerGroup,
		eventProcessor,
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	go kafkaConsumer.Start(ctx)
	slog.Info("Kafka consumer started", "topic", cfg.Kafka.Topic, "group", cfg.Kafka.ConsumerGroup)

	defer func() {
		if err := kafkaConsumer.Stop(); err != nil {
			slog.Error("Failed to stop Kafka consumer", "error", err)
		}
	}()
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.GlobalRPS, cfg.RateLimit.IPRPS)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.RateLimitMiddleware(rateLimiter))
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Metrics)

	r.Handle("/metrics", promhttp.Handler())

	handler.SetupRoutes(r, courierService, deliveryService)

	go func() {
		pprofRouter := http.NewServeMux()
		
		pprofRouter.HandleFunc("/debug/pprof/", pprof.Index)
		pprofRouter.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofRouter.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofRouter.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		pprofRouter.HandleFunc("/debug/pprof/trace", pprof.Trace)
		
		pprofRouter.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		pprofRouter.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		pprofRouter.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		pprofRouter.Handle("/debug/pprof/block", pprof.Handler("block"))
		pprofRouter.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
		pprofRouter.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
		
		slog.Info("Starting pprof server", "addr", "0.0.0.0:6060")
		if err := http.ListenAndServe("0.0.0.0:6060", pprofRouter); err != nil {
			slog.Error("pprof server failed", "error", err)
		}
	}()

	port := strconv.Itoa(cfg.Port)
	slog.Info("Starting server", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func setupSignalHandler(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		slog.Info("Received signal, shutting down", "signal", sig)
		cancel()

		time.Sleep(5 * time.Second)
		slog.Info("Shutdown complete")
		os.Exit(0)
	}()
}