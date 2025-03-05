package boot

import (
	"distributed-lock/config"
	"distributed-lock/controller"
	"distributed-lock/database"
	"distributed-lock/database/invoice"
	"distributed-lock/database/order"
	"distributed-lock/locker"
	"distributed-lock/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Order(cfg config.Config) (*gin.Engine, error) {
	db, err := database.OpenDB(config.RelationalDatabase{
		Host:     cfg.OrderDatabase.Host,
		Port:     cfg.OrderDatabase.Port,
		User:     cfg.OrderDatabase.User,
		Password: cfg.OrderDatabase.Password,
		Name:     cfg.OrderDatabase.Name,
	}, order.MigrationsFS)
	if err != nil {
		return nil, fmt.Errorf("failed open database: %w", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.NonRelationalDatabase.Host, cfg.NonRelationalDatabase.Port),
	})
	lockManager := locker.NewLockManager(*client)
	repo := repository.NewOrderDB(db)
	httpClient := http.DefaultClient
	controller := controller.NewOrder(&repo, httpClient, lockManager)
	return newRouter(controller), nil
}

func Invoice(cfg config.Config) (*gin.Engine, error) {
	db, err := database.OpenDB(config.RelationalDatabase{
		Host:     cfg.InvoiceDatabase.Host,
		Port:     cfg.InvoiceDatabase.Port,
		User:     cfg.InvoiceDatabase.User,
		Password: cfg.InvoiceDatabase.Password,
		Name:     cfg.InvoiceDatabase.Name,
	}, invoice.MigrationsFS)
	if err != nil {
		return nil, fmt.Errorf("failed open database: %w", err)
	}
	repo := repository.NewInvoiceDB(db)
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.NonRelationalDatabase.Host, cfg.NonRelationalDatabase.Port),
	})
	lockManager := locker.NewLockManager(*client)
	return newRouterInvoice(controller.NewInvoice(repo, lockManager)), nil
}
