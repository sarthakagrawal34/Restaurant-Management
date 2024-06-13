package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management/database"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var OrderCollection *mongo.Collection = database.OpenCollection(database.MongoClient, "order")

func GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := OrderCollection.Find(ctx, bson.M{})
		if err != nil {
			fmt.Println("error in GetAllOrders function while finding order items, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items"})
			return
		}

		var allOrders []models.Order
		if err = result.All(ctx, &allOrders); err != nil {
			fmt.Println("error in GetAllOrders function while decoding all orders, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order details"})
			return
		}
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderId := c.Param("order_id")
		var order models.Order
		err := OrderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		if err != nil {
			fmt.Println("error in GetOrder function in finding order, err: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching order details"})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			fmt.Println("error in CreateOrder function while binding order json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding order object"})
			return
		}

		validationError := validate.Struct(order)
		if validationError != nil {
			fmt.Println("error in CreateOrder function while validating Order json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating Order object"})
			return
		}

		var table models.RestaurantTable
		err := restaurantTableCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
		if err != nil {
			fmt.Println("error in CreateOrder function while finding linked restaturant table, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "restaturant table not found"})
			return
		}

		order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.OrderId = order.ID.Hex()

		res, err := OrderCollection.InsertOne(ctx, order)
		if err != nil {
			fmt.Println("error in CreateOrder function while creating order, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in creating order"})
			return
		}

		msg := fmt.Sprintf("Order is inserted successfully with the insertion number as: %v", res)
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			fmt.Println("error in UpdateOrder function while binding order json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding order object"})
			return
		}

		validationError := validate.Struct(order)
		if validationError != nil {
			fmt.Println("error in UpdateOrder function while validating order json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating order object"})
			return
		}

		var updateObj primitive.D

		if order.TableId != nil {
			var table models.RestaurantTable
			err := restaurantTableCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
			if err != nil {
				fmt.Println("error in UpdateOrder function while finding linked table, err: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "table not found"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.TableId})
		}

		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		orderId := c.Param("order_id")
		filter := bson.M{"order_id": orderId}
		res, err := OrderCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
		if err != nil {
			fmt.Println("error in UpdateOrder function while updating order, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while updating order"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

