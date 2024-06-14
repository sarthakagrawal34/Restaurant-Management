package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management/helpers"
	"restaurant-management/models"
	"restaurant-management/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := models.OrderItemCollection.Find(ctx, bson.M{})
		if err != nil {
			fmt.Println("error in GetAllOrderItems function while finding order items, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items"})
			return
		}

		var allOrderItems []models.OrderItem
		if err = result.All(ctx, &allOrderItems); err != nil {
			fmt.Println("error in GetAllOrderItems function while decoding all orders items, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items"})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderItemId := c.Param("order_item_id")
		var order_item models.OrderItem
		err := models.OrderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&order_item)
		if err != nil {
			fmt.Println("error in GetOrderItem function in finding order item, err: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching order item details"})
			return
		}
		c.JSON(http.StatusOK, order_item)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderId := c.Param("order_id")
		allOrderItems, err := helpers.ItemsByOrder(ctx, orderId)

		if err != nil {
			fmt.Printf("error in GetOrderItemsByOrder function while finding all orders items for a orderId: %v, err: %v\n", orderId, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items"})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)

	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var orderItemPack models.OrderItemPack
		if err := c.BindJSON(&orderItemPack); err != nil {
			fmt.Println("error in CreateOrderItem function while binding orderItemPack json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding orderItemPack object"})
			return
		}

		var order models.Order
		order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemsToBeInserted := []interface{}{}
		order.TableId = orderItemPack.TableId
		order_id, err := helpers.OrderItemOrderCreator(ctx, order)
		if err != nil {
			fmt.Println("error in CreateOrderItem function while binding orderItemPack json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding orderItemPack object"})
			return
		}

		for _, orderItem := range orderItemPack.OrderItems {
			orderItem.OrderId = &order_id

			validationError := validate.Struct(orderItem)
			if validationError != nil {
				fmt.Println("error in CreateOrderItem function while validating orderItem json, err: ", validationError.Error())
				c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating orderItem object"})
				return
			}

			orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.ID = primitive.NewObjectID()
			orderItem.OrderItemId = orderItem.ID.Hex()
			var num = utils.ToFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num

			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		res, err := models.OrderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			fmt.Println("error in CreateOrderItem function while creating OrderItem, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in creating OrderItem"})
			return
		}

		msg := fmt.Sprintf("OrderItems are inserted successfully with the insertion numbers as: %v", res)
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var orderItem models.OrderItem
		if err := c.BindJSON(&orderItem); err != nil {
			fmt.Println("error in UpdateOrderItem function while binding orderItem json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding orderItem object"})
			return
		}

		validationError := validate.Struct(orderItem)
		if validationError != nil {
			fmt.Println("error in UpdateOrderItem function while validating orderItem json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating orderItem object"})
			return
		}

		var updateObj primitive.D

		if orderItem.UnitPrice != nil {
			var num = utils.ToFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num
			updateObj = append(updateObj, bson.E{Key: "unit_price", Value: orderItem.UnitPrice})
		}
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
		}
		if orderItem.FoodId != nil {
			updateObj = append(updateObj, bson.E{Key: "food_id", Value: orderItem.FoodId})
		}
		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		orderItemId := c.Param("order_item_id")
		filter := bson.M{"order_item_id": orderItemId}
		res, err := models.OrderItemCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
		if err != nil {
			fmt.Println("error in UpdateOrderItem function while updating orderItem, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while updating orderItem"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
