package controller

import (
	"distributed-lock/invoice/postgres"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c Invoice) CreateInvoice(ctx *gin.Context) {
	var body CreateInvoiceInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, output{Err: err.Error(), Message: "error converting request body"})
		return
	}

	invoice := postgres.InvoiceModel{
		ID:        uuid.NewString(),
		UserID:    body.UserID,
		Status:    postgres.InvoiceStatusCreated,
		Total:     body.Total,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := c.repo.Create(ctx, invoice); err != nil {
		ctx.JSON(http.StatusInternalServerError, output{Err: err.Error(), Message: "internal server error"})
		return
	}
	ctx.JSON(http.StatusCreated, CreateInvoiceOutput{ID: invoice.ID})
}

type CreateInvoiceInput struct {
	UserID string  `json:"user_id"`
	Total  float64 `json:"total"`
}

type CreateInvoiceOutput struct {
	ID string `json:"invoice_id"`
}
