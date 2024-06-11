package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func RestaurantTableRoutes(in *gin.Engine) {
	in.GET("/tables", controller.GetAllRestaurantTables())
	in.GET("/tables/:table_id", controller.GetRestaurantTable())
	in.POST("/tables", controller.CreateRestaurantTable())
	in.PATCH("/tables/:table_id", controller.UpdateRestaurantTable())
}