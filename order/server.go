package order

import (
	"distributed-lock/order/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func newRouter(controller controller.Order) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})
	r.POST("/orders", controller.CreateOrder)
	r.DELETE("/orders/:id", controller.DeleteOrder)
	r.PATCH("/orders/:id", controller.UpdateOrder)
	return r
}

func newRouterInvoice(c controller.Invoice) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})
	r.POST("/invoices", c.CreateInvoice)
	r.DELETE("/invoices/:id", c.DeleteInvoice)
	r.PATCH("/invoices/:id", c.UpdateInvoice)
	return r
}
