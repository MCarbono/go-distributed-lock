package controller

import (
	"distributed-lock/model"
	"distributed-lock/repository"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c Invoice) UpdateInvoice(ctx *gin.Context) {
	id := ctx.Param("id")
	var body UpdateInvoiceInput
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, output{Err: err.Error(), Message: "error converting request body"})
		return
	}
	locked, err := c.locker.ExistsLock(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, output{Err: err.Error()})
		return
	}
	if !locked {
		ctx.JSON(http.StatusConflict, output{Message: fmt.Sprintf("resource with id %s should be locked but its not.", id)})
		return
	}
	invoice, err := c.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrInvoiceNotFound) {
			ctx.JSON(http.StatusNotFound, output{Err: err.Error(), Message: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, output{Err: err.Error(), Message: "error trying to get the invoice"})
		return
	}
	invoice.UpdatedAt = time.Now().UTC()
	invoice.Status = model.InvoiceStatusUpdated
	invoice.Total = body.Total
	err = c.repo.Update(ctx, invoice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, output{Err: err.Error(), Message: "internal server error"})
		return
	}
	time.Sleep(2 * time.Second)
	ctx.JSON(http.StatusNoContent, nil)
}
