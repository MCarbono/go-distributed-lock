package boot

import (
	"distributed-lock/controller"
	invoiceDatabase "distributed-lock/database/invoice"
	orderDatabase "distributed-lock/database/order"
	"distributed-lock/locker"
	"distributed-lock/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Order() (*gin.Engine, error) {
	db, err := orderDatabase.OpenDB()
	if err != nil {
		return nil, fmt.Errorf("failed open database: %w", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	lockManager := locker.NewLockManager(*client)
	repo := repository.NewOrderDB(db)
	httpClient := http.DefaultClient
	controller := controller.NewOrder(&repo, httpClient, lockManager)
	return newRouter(controller), nil
}

func Invoice() (*gin.Engine, error) {
	db, err := invoiceDatabase.OpenDBInvoice()
	if err != nil {
		return nil, fmt.Errorf("failed open database: %w", err)
	}
	repo := repository.NewInvoiceDB(db)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	lockManager := locker.NewLockManager(*client)
	return newRouterInvoice(controller.NewInvoice(repo, lockManager)), nil
}
