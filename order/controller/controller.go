package controller

import (
	"context"
	"distributed-lock/locker"
	"distributed-lock/order/postgres"
	"distributed-lock/order/repository"
	"fmt"
	"net/http"
)

type Order struct {
	repo        repository.Order
	httpClient  *http.Client
	lockManager locker.LockManager
}

func NewOrder(repo repository.Order, client *http.Client, lockerManager locker.LockManager) Order {
	return Order{
		repo:        repo,
		httpClient:  client,
		lockManager: lockerManager,
	}
}

type errOutput struct {
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (c Order) releaseLocks(ctx context.Context, order postgres.OrderModel) {
	if err := c.lockManager.ReleaseLock(ctx, order.ID); err != nil {
		fmt.Printf("unlocking order failed: %v\n", err)
	}
	if err := c.lockManager.ReleaseLock(ctx, order.InvoiceID); err != nil {
		fmt.Printf("unlocking invoice failed: %v\n", err)
	}
}
