package order

import (
	"distributed-lock/locker"
	"distributed-lock/order/controller"
	"distributed-lock/order/postgres"
	"distributed-lock/order/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Boot() (*gin.Engine, error) {
	db, err := postgres.OpenDB()
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
