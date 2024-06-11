package main

import (
	"os"
	"restaurant-management/database"
	"restaurant-management/middleware"
	"restaurant-management/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.MongoClient, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	// define a new engine 
	router := gin.New()
	// Add logger library for rich logs
	router.Use(gin.Logger())
	// use user routes for login and signup
	routes.UserRoutes(router)
	// pass the routes through middleware for authentication
	router.Use(middleware.Authenticate())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.RestaurantTableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	// run the server on port: 8080
	router.Run(":" + port)
}