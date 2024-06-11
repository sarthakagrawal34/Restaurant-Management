package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(in *gin.Engine) {
	in.GET("/orders", controller.GetAllOrders())
	in.GET("/orders/:order_id", controller.GetOrder())
	in.POST("/orders", controller.CreateOrder())
	in.PATCH("/orders/:order_id", controller.UpdateOrder())
}