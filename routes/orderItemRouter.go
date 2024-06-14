package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(in *gin.Engine) {
	in.GET("/orderItems", controller.GetAllOrderItems())
	in.GET("/orderItems/:order_item_id", controller.GetOrderItem())
	in.GET("/orderItems-order/:order_id", controller.GetOrderItemsByOrder())
	in.POST("/orderItems", controller.CreateOrderItem())
	in.PATCH("/orderItems/:order_item_id", controller.UpdateOrderItem())
}