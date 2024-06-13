package helpers

import (
	"context"
	"fmt"
	"restaurant-management/controllers"
	"restaurant-management/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	return nil, nil
}

func OrderItemOrderCreator(order models.Order) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderId = order.ID.Hex()

	_, err := controllers.OrderCollection.InsertOne(ctx, order)
	if err != nil {
		fmt.Println("error in CreateOrder function while creating order, err: ", err)
		return "", err
	}
	return order.OrderId, nil
}
