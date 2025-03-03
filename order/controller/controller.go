package controller

import (
	"distributed-lock/locker"
	"distributed-lock/order/repository"
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
