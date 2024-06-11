package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(in *gin.Engine) {
	in.GET("/foods", controller.GetAllFoods())
	in.GET("/foods/:food_id", controller.GetFood())
	in.POST("/foods", controller.CreateFood())
	in.PATCH("/foods/:food_id", controller.UpdateFood())
}
