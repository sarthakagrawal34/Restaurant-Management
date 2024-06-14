package routes

import (
	controller "restaurant-management/controllers"
	"restaurant-management/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(in *gin.Engine) {
	in.POST("/users/signup", controller.SignUp())
	in.POST("/users/login", controller.Login())
	in.Use(middleware.Authenticate())
	in.GET("/users", controller.GetAllUsers())
	in.GET("/users/:user_id", controller.GetUser())
}