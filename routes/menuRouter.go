package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(in *gin.Engine) {
	in.GET("/menus", controller.GetAllMenus())
	in.GET("/menus/:menu_id", controller.GetMenu())
	in.POST("/menus", controller.CreateMenu())
	in.PATCH("/menus/:menu_id", controller.UpdateMenu())
}