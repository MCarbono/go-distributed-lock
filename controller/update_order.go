package controller

import (
	"bytes"
	"distributed-lock/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c Order) UpdateOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	var input UpdateOrderInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errOutput{Err: err.Error(), Message: "error converting request body"})
		return
	}
	if input.UpdateTime.IsZero() {
		err := errors.New("update_time is a required field")
		ctx.JSON(http.StatusUnprocessableEntity, errOutput{Err: err.Error(), Message: err.Error()})
		return
	}
	order, err := c.repo.FindByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error()})
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
	if order.Status == model.OrderStatusDeleted {
		ctx.JSON(http.StatusUnprocessableEntity, errOutput{Err: "order_cancelled", Message: "order is cancelled already"})
		return
	}
	if input.UpdateTime.UTC().Before(order.UpdatedAt.UTC()) {
		err := fmt.Errorf("the update order request is outdated. Current instant: %s, last updated instant: %s", input.UpdateTime.UTC(), order.UpdatedAt.UTC())
		ctx.JSON(http.StatusConflict, errOutput{Err: "event_outdated", Message: err.Error()})
		return
	}
	if input.hasQuantity() {
		order.Quantity = input.Quantity
	}
	if input.hasValue() {
		order.Value = input.Value
	}
	order.Total = float64(order.Quantity) * order.Value
	invoiceInput, err := json.Marshal(UpdateInvoiceInput{Total: order.Total})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error(), Message: "error converting invoice request body"})
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("http://localhost:3001/invoices/%s", order.InvoiceID), bytes.NewBuffer(invoiceInput))
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
	order.UpdatedAt = input.UpdateTime.UTC()
	order.Status = model.OrderStatusUpdated
	err = c.repo.Update(ctx, order)
	if err != nil {
		//TODO - undo update invoice
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error(), Message: "error updating order into the database"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

type UpdateOrderInput struct {
	Quantity   int       `json:"quantity"`
	Value      float64   `json:"value"`
	UpdateTime time.Time `json:"update_time"`
}

func (u UpdateOrderInput) hasValue() bool {
	return u.Value > 0
}

func (u UpdateOrderInput) hasQuantity() bool {
	return u.Quantity > 0
}

type UpdateInvoiceInput struct {
	Total float64 `json:"total"`
}
