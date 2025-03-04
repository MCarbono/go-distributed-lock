package controller

import (
	"distributed-lock/locker"
	"distributed-lock/repository"
)

type Invoice struct {
	repo   repository.Invoice
	locker locker.LockManager
}

func NewInvoice(repo repository.Invoice, locker locker.LockManager) Invoice {
	return Invoice{
		repo:   repo,
		locker: locker,
	}
}

type output struct {
	Err     string `json:"error"`
	Message string `json:"message"`
}
