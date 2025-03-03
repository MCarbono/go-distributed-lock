package invoice

import (
	"distributed-lock/invoice/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func newRouter(c controller.Invoice) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})
	r.POST("/invoices", c.CreateInvoice)
	r.DELETE("/invoices/:id", c.DeleteInvoice)
	r.PATCH("/invoices/:id", c.UpdateInvoice)
	return r
}
