package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management/helpers"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := models.InvoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			fmt.Println("error in GetAllInvoices function while finding invoice items, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing invoice items"})
			return
		}

		var allInvoices []models.Invoice
		if err = result.All(ctx, &allInvoices); err != nil {
			fmt.Println("error in GetAllInvoices function while decoding all invoices, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing invoice details"})
			return
		}
		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		invoiceId := c.Param("invoice_id")

		var invoice models.Invoice
		err := models.InvoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		if err != nil {
			fmt.Println("error in GetInvoice function in finding invoice, err: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching invoice details"})
			return
		}

		var invoiceView models.InvoiceViewFormat
		allOrderItems, err := helpers.ItemsByOrder(*invoice.OrderId)
		if err != nil {
			fmt.Println("error in GetInvoice function in finding all order item details, err: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching invoice order items details"})
			return
		}
		invoiceView.OrderId = *invoice.OrderId
		invoiceView.PaymentDueDate = invoice.PaymentDueDate
		invoiceView.PaymentMethod = "null"
		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = *invoice.PaymentMethod
		}
		invoiceView.InvoiceId = invoice.InvoiceId
		invoiceView.PaymentStatus = invoice.PaymentStatus
		invoiceView.PaymentDue = allOrderItems[0]["payment_due"]
		invoiceView.TableNumber = allOrderItems[0]["table_number"]
		invoiceView.OrderDetails = allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK, invoiceView)

	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			fmt.Println("error in CreateInvoice function while binding invoice json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding invoice object"})
			return
		}

		validationError := validate.Struct(invoice)
		if validationError != nil {
			fmt.Println("error in CreateInvoice function while validating invoice json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating invoice object"})
			return
		}

		var order models.Order
		err := models.OrderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderId}).Decode(&order)
		if err != nil {
			fmt.Println("error in CreateInvoice function while finding linked order, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order not found"})
			return
		}

		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}
		invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceId = invoice.ID.Hex()

		res, err := models.InvoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			fmt.Println("error in CreateInvoice function while creating invoice, err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in creating invoice"})
			return
		}

		msg := fmt.Sprintf("Invoice is inserted successfully with the insertion number as: %v", res)
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			fmt.Println("error in UpdateInvoice function while binding invoice json, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while binding invoice object"})
			return
		}

		validationError := validate.Struct(invoice)
		if validationError != nil {
			fmt.Println("error in UpdateInvoice function while validating invoice json, err: ", validationError.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while validating invoice object"})
			return
		}

		var updateObj primitive.D

		if invoice.PaymentMethod != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_method", Value: invoice.PaymentMethod})
		}

		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}
		updateObj = append(updateObj, bson.E{Key: "payment_status", Value: invoice.PaymentStatus})

		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		invoiceId := c.Param("invoice_id")
		filter := bson.M{"invoice_id": invoiceId}
		res, err := models.InvoiceCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
		if err != nil {
			fmt.Println("error in UpdateInvoice function while updating invoice, err: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while updating invoice"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
