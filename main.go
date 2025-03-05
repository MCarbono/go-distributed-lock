package main

import (
	"context"
	"distributed-lock/boot"
	"distributed-lock/config"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	env = flag.String("env", "local", "used to know what environment the project is running")
)

func main() {
	flag.Parse()
	cfg, err := config.LoadEnv(*env)
	if err != nil {
		panic(cfg)
	}
	fmt.Printf("Env variables loaded: %+v\n", cfg)

	orderService := http.Server{Addr: fmt.Sprintf(":%s", cfg.OrderServerPort)}
	invoiceService := http.Server{Addr: fmt.Sprintf(":%s", cfg.InvoiceServerPort)}

	ch := make(chan error)
	defer close(ch)
	go newOrderServiceRouter(ch, &orderService, cfg)
	err = <-ch
	if err != nil {
		panic(err)
	}

	go newInvoiceServiceRouter(ch, &invoiceService, cfg)
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

func newOrderServiceRouter(ch chan<- error, service *http.Server, cfg config.Config) {
	router, err := boot.Order(cfg.OrderDatabase, cfg.NonRelationalDatabase)
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

func newInvoiceServiceRouter(ch chan<- error, service *http.Server, cfg config.Config) {
	router, err := boot.Invoice(cfg.InvoiceDatabase, cfg.NonRelationalDatabase)
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
