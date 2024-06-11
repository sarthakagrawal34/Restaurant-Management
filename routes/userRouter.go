package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(in *gin.Engine) {
	in.GET("/users", controller.GetAllUsers())
	in.GET("/users/:user_id", controller.GetUser())
	in.POST("/users/signup", controller.SignUp())
	in.POST("/users/login", controller.Login())
}