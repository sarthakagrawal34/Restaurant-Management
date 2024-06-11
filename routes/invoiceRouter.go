package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(in *gin.Engine) {
	in.GET("/invoices", controller.GetAllInvoices())
	in.GET("/invoices/:invoice_id", controller.GetInvoice())
	in.POST("/invoices", controller.CreateInvoice())
	in.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())
}