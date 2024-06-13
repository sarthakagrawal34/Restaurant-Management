package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := models.RestaurantTableCollection.Find(ctx, bson.M{})
		if err != nil {
			fmt.Println("error in GetAllTables function while finding all tables, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing all tables"})
			return
		}

		var allTables []models.RestaurantTable
		if err = result.All(ctx, &allTables); err != nil {
			fmt.Println("error in GetAllTables function while decoding all tables, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing all tables"})
			return
		}
		c.JSON(http.StatusOK, allTables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tableId := c.Param("table_id")
		var table models.RestaurantTable
		err := models.RestaurantTableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			fmt.Println("error in GetTables function in finding table, err: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the table details"})
			return
		}
		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var table models.RestaurantTable
		if err := c.BindJSON(&table); err != nil {
			fmt.Println("error in CreateTable function while binding table json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding table object"})
			return
		}

		validationError := validate.Struct(table)
		if validationError != nil {
			fmt.Println("error in CreateTable function while validating table json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating table object"})
			return
		}

		table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.TableId = table.ID.Hex()

		res, err := models.RestaurantTableCollection.InsertOne(ctx, table)
		if err != nil {
			fmt.Println("error in CreateTable function while creating table, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in creating table"})
			return
		}

		msg := fmt.Sprintf("Table is inserted successfully with the insertion number as: %v", res)
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var table models.RestaurantTable
		if err := c.BindJSON(&table); err != nil {
			fmt.Println("error in UpdateTable function while binding table json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding table object"})
			return
		}

		validationError := validate.Struct(table)
		if validationError != nil {
			fmt.Println("error in UpdateTable function while validating table json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating table object"})
			return
		}

		var updateObj primitive.D

		if table.NumberOfGuests != nil {
			updateObj = append(updateObj, bson.E{Key: "number_of_guests", Value: table.NumberOfGuests})
		}
		if table.TableNumber != nil {
			updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.TableNumber})
		}
		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		tableId := c.Param("table_id")
		filter := bson.M{"table_id": tableId}
		res, err := models.RestaurantTableCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
		if err != nil {
			fmt.Println("error in UpdateTable function while updating table, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while updating table"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
