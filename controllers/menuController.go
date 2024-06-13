package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management/database"
	"restaurant-management/models"
	"restaurant-management/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.MongoClient, "menu")

func GetAllMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := menuCollection.Find(ctx, bson.M{})
		if err != nil {
			fmt.Println("error in GetAllMenus function while finding menu items, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing menu items"})
			return
		}

		var allMenus []models.Menu
		if err = result.All(ctx, &allMenus); err != nil {
			fmt.Println("error in GetAllMenus function while decoding all menus, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing menu details"})
			return
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			fmt.Println("error in GetMenu function in finding menu, err: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching menu details"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			fmt.Println("error in CreateMenu function while binding menu json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding menu object"})
			return
		}

		validationError := validate.Struct(menu)
		if validationError != nil {
			fmt.Println("error in CreateMenu function while validating menu json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating menu object"})
			return
		}

		if menu.StartDate != nil && menu.EndDate != nil {
			if !utils.InTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
				fmt.Println("error in CreateMenu function while validating start and end date of menu")
				msg := "kindly retype the valid time"
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				return
			}
		}

		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.MenuId = menu.ID.Hex()

		res, err := menuCollection.InsertOne(ctx, menu)
		if err != nil {
			fmt.Println("error in CreateMenu function while creating menu, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in creating menu"})
			return
		}

		msg := fmt.Sprintf("Menu is inserted successfully with the insertion number as: %v", res)
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			fmt.Println("error in UpdateMenu function while binding menu json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding menu object"})
			return
		}

		validationError := validate.Struct(menu)
		if validationError != nil {
			fmt.Println("error in UpdateMenu function while validating menu json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating menu object"})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj primitive.D
		if menu.StartDate != nil && menu.EndDate != nil {
			if !utils.InTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
				fmt.Println("error in UpdateMenu function while validating start and end date of menu")
				msg := "kindly retype the valid time"
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.StartDate})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.EndDate})

			if menu.Name != nil {
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
			}
			if menu.Category != nil {
				updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
			}
			menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.UpdatedAt})

			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}

			res, err := menuCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
			if err != nil {
				fmt.Println("error in UpdateMenu function while updating menu, err: ", err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while updating menu"})
				return
			}
			c.JSON(http.StatusOK, res)
		}
	}
}
