package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

type RateLimitConfig struct {
	GlobalRPS int
	IPRPS     int
}

type GatewayRetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
}

type Config struct {
	Port            int
	DB              DatabaseConfig
	OrderServiceURL string
	Kafka           KafkaConfig
	RateLimit       RateLimitConfig
	GatewayRetry    GatewayRetryConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type KafkaConfig struct {
	Broker        string
	Topic         string
	ConsumerGroup string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	portFlag := pflag.Int("port", 0, "HTTP server port")
	pflag.Parse()

	// По умолчанию
	port := 8080
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}
	if *portFlag != 0 {
		port = *portFlag
	}

	// По умолчанию
	dbPort := 5432
	if envDBPort := os.Getenv("POSTGRES_PORT"); envDBPort != "" {
		if p, err := strconv.Atoi(envDBPort); err == nil {
			dbPort = p
		}
	}

	orderServiceURL := getEnv("ORDER_SERVICE_HOST", "service-order:50051")

	return &Config{
		Port: port,
		DB: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("POSTGRES_USER", "myuser"),
			Password: getEnv("POSTGRES_PASSWORD", "mypassword"),
			DBName:   getEnv("POSTGRES_DB", "test_db"),
		},
		OrderServiceURL: orderServiceURL,
		Kafka: KafkaConfig{
			Broker:        getEnv("KAFKA_BROKER", "localhost:9093"),
			Topic:         getEnv("KAFKA_TOPIC", "order.status.changed"),
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "service-courier-group"),
		},
		RateLimit: RateLimitConfig{
			GlobalRPS: getEnvAsInt("RATE_LIMIT_GLOBAL_RPS", 100),
			IPRPS:     getEnvAsInt("RATE_LIMIT_IP_RPS", 10),
		},
		GatewayRetry: GatewayRetryConfig{
			MaxAttempts:  getEnvAsInt("GATEWAY_RETRY_MAX_ATTEMPTS", 3),
			InitialDelay: getEnvAsDuration("GATEWAY_RETRY_INITIAL_DELAY_MS", 100) * time.Millisecond,
			MaxDelay:     getEnvAsDuration("GATEWAY_RETRY_MAX_DELAY_MS", 2000) * time.Millisecond,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue int) time.Duration {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return time.Duration(intValue)
		}
	}
	return time.Duration(defaultValue)
}
