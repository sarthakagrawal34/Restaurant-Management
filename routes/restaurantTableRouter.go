package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(in *gin.Engine) {
	in.GET("/tables", controller.GetAllTables())
	in.GET("/tables/:table_id", controller.GetTable())
	in.POST("/tables", controller.CreateTable())
	in.PATCH("/tables/:table_id", controller.UpdateTable())
}
