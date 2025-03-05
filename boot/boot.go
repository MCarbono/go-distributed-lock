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

func Order(relationalDatabaseConfing config.RelationalDatabaseOrderService, nonRelationalDatabase config.NonRelationalDatabase) (*gin.Engine, error) {
	db, err := database.OpenDB(config.RelationalDatabase{
		Host:     relationalDatabaseConfing.Host,
		Port:     relationalDatabaseConfing.Port,
		User:     relationalDatabaseConfing.User,
		Password: relationalDatabaseConfing.Password,
		Name:     relationalDatabaseConfing.Name,
	}, order.MigrationsFS)
	if err != nil {
		return nil, fmt.Errorf("failed open database: %w", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", nonRelationalDatabase.Host, nonRelationalDatabase.Port),
	})
	lockManager := locker.NewLockManager(*client)
	repo := repository.NewOrderDB(db)
	httpClient := http.DefaultClient
	controller := controller.NewOrder(&repo, httpClient, lockManager)
	return newRouter(controller), nil
}

func Invoice(relationalDatabaseConfing config.RelationalDatabaseInvoiceService, nonRelationalDatabase config.NonRelationalDatabase) (*gin.Engine, error) {
	db, err := database.OpenDB(config.RelationalDatabase{
		Host:     relationalDatabaseConfing.Host,
		Port:     relationalDatabaseConfing.Port,
		User:     relationalDatabaseConfing.User,
		Password: relationalDatabaseConfing.Password,
		Name:     relationalDatabaseConfing.Name,
	}, invoice.MigrationsFS)
	if err != nil {
		return nil, fmt.Errorf("failed open database: %w", err)
	}
	repo := repository.NewInvoiceDB(db)
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", nonRelationalDatabase.Host, nonRelationalDatabase.Port),
	})
	lockManager := locker.NewLockManager(*client)
	return newRouterInvoice(controller.NewInvoice(repo, lockManager)), nil
}
