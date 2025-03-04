package main

import (
	"context"
	"distributed-lock/boot"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	orderService := http.Server{Addr: ":3000"}
	invoiceService := http.Server{Addr: ":3001"}

	ch := make(chan error)
	defer close(ch)
	go newOrderServiceRouter(ch, &orderService)
	err := <-ch
	if err != nil {
		panic(err)
	}

	go newInvoiceServiceRouter(ch, &invoiceService)
	err = <-ch
	if err != nil {
		shutDownCtx, shutdowRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdowRelease()
		if err := orderService.Shutdown(shutDownCtx); err != nil {
			fmt.Printf("orderService shutdown error: %v", err.Error())
		}
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutDownCtx, shutdowRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdowRelease()
	if err := orderService.Shutdown(shutDownCtx); err != nil {
		fmt.Printf("orderService shutdown error: %v", err.Error())
	}
	fmt.Println("orderService terminated successfully")
	if err := invoiceService.Shutdown(shutDownCtx); err != nil {
		fmt.Printf("invoiceService shutdown error: %v", err.Error())
	}
	fmt.Println("invoiceService terminated successfully")
}

func newOrderServiceRouter(ch chan<- error, service *http.Server) {
	router, err := boot.Order()
	if err != nil {
		ch <- err
		return
	}
	ch <- nil
	service.Handler = router
	if err := service.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("stopped serving connection for orderService")
			return
		}
		fmt.Printf("orderService shutdown error: %v", err)
	}
}

func newInvoiceServiceRouter(ch chan<- error, service *http.Server) {
	router, err := boot.Invoice()
	if err != nil {
		ch <- err
		return
	}
	ch <- nil
	service.Handler = router
	if err := service.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("stopped serving connection for invoiceService")
			return
		}
		fmt.Printf("invoiceService shutdown error: %v", err)
	}
}
