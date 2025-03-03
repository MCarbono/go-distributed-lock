package controller

import (
	"database/sql"
	"distributed-lock/invoice/postgres"
	"distributed-lock/invoice/repository"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c Invoice) DeleteInvoice(ctx *gin.Context) {
	id := ctx.Param("id")
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
	if invoice.Status == postgres.InvoiceStatusDeleted {
		ctx.JSON(http.StatusNoContent, nil)
		return
	}
	invoice.UpdatedAt = time.Now().UTC()
	invoice.DeletedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	invoice.Status = postgres.InvoiceStatusDeleted
	err = c.repo.Delete(ctx, invoice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, output{Err: err.Error(), Message: "internal server error"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
