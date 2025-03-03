package invoice

import (
	"distributed-lock/invoice/controller"
	"distributed-lock/invoice/postgres"
	"distributed-lock/invoice/repository"
	"distributed-lock/locker"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Boot() (*gin.Engine, error) {
	db, err := postgres.OpenDBInvoice()
	if err != nil {
		return nil, fmt.Errorf("failed open database: %w", err)
	}
	repo := repository.NewInvoiceDB(db)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	lockManager := locker.NewLockManager(*client)
	return newRouter(controller.NewInvoice(repo, lockManager)), nil
}
