package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"My_WEB/config"
	"My_WEB/models"
	"My_WEB/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateProduct
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Lấy collection từ config.DB
	productCollection := config.DB.Collection("products")
	if productCollection == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database collection not initialized")
		return
	}

	// Parse dữ liệu từ request
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Thêm ID tự động
	product.ID = primitive.NewObjectID()

	// Chèn vào MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to insert product")
		return
	}

	// Log để debug
	fmt.Printf("Product inserted with ID: %v\n", result.InsertedID)

	// Phản hồi thành công với ID của sản phẩm vừa tạo
	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Product created successfully",
		"id":      result.InsertedID,
	})
}

// GetAllProducts
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productCollection := config.DB.Collection("products")
	if productCollection == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database collection not initialized")
		return
	}

	var products []models.Product

	// Tạo context với timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tìm tất cả tài liệu trong collection "products"
	cursor, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer cursor.Close(ctx)

	// Lặp qua từng tài liệu trong cursor và giải mã vào slice products
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			fmt.Printf("Error decoding product: %v\n", err)
			continue // Bỏ qua document lỗi và tiếp tục
		}
		products = append(products, product)
	}

	// Kiểm tra lỗi nếu có trong quá trình duyệt cursor
	if err := cursor.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log để debug
	fmt.Printf("Found %d products\n", len(products))

	// Trả về danh sách sản phẩm dưới dạng JSON với mã 200 OK
	utils.RespondWithJSON(w, http.StatusOK, products)
}

// GetProductByID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productCollection := config.DB.Collection("products")
	if productCollection == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database collection not initialized")
		return
	}

	params := mux.Vars(r)
	productID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	var product models.Product
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Log để debug
	fmt.Printf("Searching for product with ID: %s\n", productID.Hex())

	// Tìm một tài liệu duy nhất có _id khớp với productID
	err = productCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log để debug
	fmt.Printf("Found product: %+v\n", product)

	// Trả về sản phẩm tìm được dưới dạng JSON với mã 200 OK
	utils.RespondWithJSON(w, http.StatusOK, product)
}

// UpdateProduct
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productCollection := config.DB.Collection("products")
	if productCollection == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database collection not initialized")
		return
	}

	params := mux.Vars(r)
	productID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tạo update document sử dụng toán tử $set để cập nhật  các trường
	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"quantity":    product.Quantity,
		},
	}

	// Cập nhật một tài liệu duy nhất khớp với productID
	result, err := productCollection.UpdateOne(ctx, bson.M{"_id": productID}, update)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Kiểm tra xem có tài liệu nào được sửa đổi hay không
	if result.ModifiedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found or no changes made")
		return
	}

	// Trả về kết quả cập nhật với mã 200 OK
	utils.RespondWithJSON(w, http.StatusOK, result)
}

// DeleteProduct xóa một sản phẩm khỏi database
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productCollection := config.DB.Collection("products")
	if productCollection == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database collection not initialized")
		return
	}

	params := mux.Vars(r)
	productID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Xóa một tài liệu duy nhất khớp với productID
	result, err := productCollection.DeleteOne(ctx, bson.M{"_id": productID})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Kiểm tra xem có tài liệu nào được xóa hay không
	if result.DeletedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Trả về thông báo thành công với mã 200 OK
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Product deleted successfully!",
	})
}
