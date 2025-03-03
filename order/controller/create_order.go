package controller

import (
	"bytes"
	"distributed-lock/order/postgres"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c Order) CreateOrder(ctx *gin.Context) {
	var body CreateOrderInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errOutput{Err: err.Error(), Message: "error converting request body"})
		return
	}
	total := float64(body.Quantity) * body.Value
	invoiceInput, err := json.Marshal(InvoiceRequestInput{UserID: body.UserID, Total: total})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error(), Message: "error converting invoice request body"})
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:3001/invoices", bytes.NewBuffer(invoiceInput))
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

	invoiceBody, err := io.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error(), Message: "error converting reading body from invoice service"})
		return
	}

	if res.StatusCode != http.StatusCreated {
		err = fmt.Errorf("request to invoice service returned an unexpected statusCode. Expected %v, got: %v. Response body: %s", http.StatusCreated, res.StatusCode, string(invoiceBody))
		ctx.JSON(http.StatusFailedDependency, errOutput{Err: err.Error(), Message: "unexpected response from invoice service"})
		return
	}
	var invoiceOutput InvoiceRequestOutput
	err = json.Unmarshal(invoiceBody, &invoiceOutput)
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, errOutput{Err: err.Error(), Message: "error converting response body to JSON from invoice service"})
		return
	}
	order := postgres.OrderModel{
		ID:        uuid.NewString(),
		UserID:    body.UserID,
		InvoiceID: invoiceOutput.ID,
		Status:    postgres.OrderStatusCreated,
		ItemID:    body.ItemID,
		Quantity:  body.Quantity,
		Value:     body.Value,
		Total:     total,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err = c.repo.Create(ctx, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errOutput{Err: err.Error(), Message: "error saving new order into the database"})
		return
	}
	ctx.JSON(http.StatusCreated, CreateOrderOutput{ID: order.ID})
}

type CreateOrderInput struct {
	UserID   string  `json:"user_id"`
	ItemID   string  `json:"item_id"`
	Quantity int     `json:"quantity"`
	Value    float64 `json:"value"`
}

type CreateOrderOutput struct {
	ID string `json:"id"`
}

type InvoiceRequestInput struct {
	UserID string  `json:"user_id"`
	Total  float64 `json:"total"`
}

type InvoiceRequestOutput struct {
	ID string `json:"invoice_id"`
}
