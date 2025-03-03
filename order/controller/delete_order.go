package controller

import (
	"context"
	"database/sql"
	"distributed-lock/order/postgres"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c Order) DeleteOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	order, err := c.repo.FindByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error()})
		return
	}
	if order.Status == postgres.OrderStatusDeleted {
		ctx.JSON(http.StatusNoContent, nil)
		return
	}
	defer c.releaseLocks(ctx, order)
	orderLock, err := c.lockManager.AcquireLock(ctx, order.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error(), Message: "could not process request. try again in a few moments"})
		return
	}
	invoiceLock, err := c.lockManager.AcquireLock(ctx, order.InvoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error(), Message: "could not process request. try again in a few moments"})
		return
	}
	if !orderLock || !invoiceLock {
		ctx.JSON(http.StatusLocked, errOutput{Message: "the resource you try to modify is currently blocked. Try again in a few moments"})
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("http://localhost:3001/invoices/%s", order.InvoiceID), nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error(), Message: "error creating the request to the invoice service"})
		return
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, errOutput{Err: err.Error(), Message: "error making the request to the invoice service"})
		return
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		invoiceBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("error converting reading body from invoice service: %v\n", err)
		}
		err = fmt.Errorf("request to invoice service returned an unexpected statusCode. Expected %v, got: %v. Response body: %s", http.StatusCreated, res.StatusCode, string(invoiceBody))
		ctx.JSON(http.StatusFailedDependency, errOutput{Err: err.Error(), Message: "unexpected response from invoice service"})
		return
	}
	order.DeletedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	order.UpdatedAt = time.Now().UTC()
	order.Status = postgres.OrderStatusDeleted
	err = c.repo.Delete(ctx, order)
	if err != nil {
		//TODO - undo cancelation on invoice service
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error(), Message: "error updating order into the database"})
		return
	}
	time.Sleep(5 * time.Second)
	ctx.JSON(http.StatusNoContent, nil)
}

func (c Order) releaseLocks(ctx context.Context, order postgres.OrderModel) {
	if err := c.lockManager.ReleaseLock(ctx, order.ID); err != nil {
		fmt.Printf("unlocking order failed: %v\n", err)
	}
	if err := c.lockManager.ReleaseLock(ctx, order.InvoiceID); err != nil {
		fmt.Printf("unlocking invoice failed: %v\n", err)
	}
}
