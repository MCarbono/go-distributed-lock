package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	OrderServerPort       string `envconfig:"ORDER_SERVER_PORT"`
	InvoiceServerPort     string `envconfig:"INVOICE_SERVER_PORT"`
	OrderDatabase         RelationalDatabaseOrderService
	InvoiceDatabase       RelationalDatabaseInvoiceService
	NonRelationalDatabase NonRelationalDatabase
}

type RelationalDatabaseOrderService struct {
	Host     string `envconfig:"ORDER_DATABASE_HOST"`
	Port     string `envconfig:"ORDER_DATABASE_PORT"`
	User     string `envconfig:"ORDER_DATABASE_USER"`
	Password string `envconfig:"ORDER_DATABASE_PASSWORD"`
	Name     string `envconfig:"ORDER_DATABASE_NAME"`
}

type RelationalDatabaseInvoiceService struct {
	Host     string `envconfig:"INVOICE_DATABASE_HOST"`
	Port     string `envconfig:"INVOICE_DATABASE_PORT"`
	User     string `envconfig:"INVOICE_DATABASE_USER"`
	Password string `envconfig:"INVOICE_DATABASE_PASSWORD"`
	Name     string `envconfig:"INVOICE_DATABASE_NAME"`
}

type RelationalDatabase struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type NonRelationalDatabase struct {
	Host string `envconfig:"NON_RELATIONAL_DATABASE_HOST"`
	Port string `envconfig:"NON_RELATIONAL_DATABASE_PORT"`
}

func LoadEnv(env string) (Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		return cfg, fmt.Errorf("failed reading .env: %w", err)
	}

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to parse env variables: %v", err)
	}

	if env == "production" {
		cfg.NonRelationalDatabase.Host = "redis"
		cfg.OrderDatabase.Host = "order-postgres"
		cfg.InvoiceDatabase.Host = "invoice-postgres"
	}

	return cfg, nil
}
