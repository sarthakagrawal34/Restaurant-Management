package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management/models"
	"restaurant-management/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()

func GetAllFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("limit"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		pageNumber, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || pageNumber < 1 {
			pageNumber = 1
		}

		startIndex := (pageNumber - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "food_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, err := models.FoodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		if err != nil {
			fmt.Println("Error in GetAllFoods function while doing aggregate, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing food items"})
			return
		}

		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			fmt.Println("Error in GetAllFoods function while decoding all foods, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing food items"})
			return
		}
		c.JSON(http.StatusOK, allFoods)
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foodId := c.Param("food_id")
		var food models.Food
		err := models.FoodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			fmt.Println("error in GetFood function in finding food, err: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching food details"})
			return
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			fmt.Println("error in CreateFood function while binding food json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding food object"})
			return
		}

		validationError := validate.Struct(food)
		if validationError != nil {
			fmt.Println("error in CreateFood function while validating food json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating food object"})
			return
		}

		var menu models.Menu
		err := models.MenuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuId}).Decode(&menu)
		if err != nil {
			fmt.Println("error in CreateFood function while finding linked menu, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu not found"})
			return
		}

		food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.FoodId = food.ID.Hex()
		var num = utils.ToFixed(*food.Price, 2)
		food.Price = &num

		res, err := models.FoodCollection.InsertOne(ctx, food)
		if err != nil {
			fmt.Println("error in CreateFood function while creating food, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in creating food"})
			return
		}

		msg := fmt.Sprintf("Food is inserted successfully with the insertion number as: %v", res)
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			fmt.Println("error in UpdateFood function while binding food json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding food object"})
			return
		}

		validationError := validate.Struct(food)
		if validationError != nil {
			fmt.Println("error in UpdateFood function while validating food json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating food object"})
			return
		}
		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
		}
		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "category", Value: food.Price})
		}
		if food.FoodImage != nil {
			updateObj = append(updateObj, bson.E{Key: "category", Value: food.FoodImage})
		}
		if food.MenuId != nil {
			var menu models.Menu
			err := models.MenuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuId}).Decode(&menu)
			if err != nil {
				fmt.Println("error in UpdateFood function while finding linked menu, err: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "menu not found"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.MenuId})
		}

		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		foodId := c.Param("food_id")
		filter := bson.M{"food_id": foodId}
		res, err := models.FoodCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
		if err != nil {
			fmt.Println("error in UpdateFood function while updating food, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while updating food"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
