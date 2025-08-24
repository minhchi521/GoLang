package controllers

import (
	"My_WEB/config"
	"My_WEB/models"
	"My_WEB/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Lấy collection từ config.DB (giống product)
	accountCollection := config.DB.Collection("accounts")
	if accountCollection == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database collection not initialized")
		return
	}

	// Parse dữ liệu từ request (giống product)
	var account models.Account
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Log để debug (giống product)
	fmt.Printf("Received account data: %+v\n", account)

	// Thêm ID tự động (giống product)
	account.Id = primitive.NewObjectID()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	// Tạo context với timeout (giống product)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Chèn vào MongoDB (giống product)
	result, err := accountCollection.InsertOne(ctx, account)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create account")
		return
	}

	// Log để debug (giống product)
	fmt.Printf("Account created with ID: %v\n", result.InsertedID)

	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Account created successfully",
		"id":      result.InsertedID,
		"user": map[string]interface{}{
			"id":       account.Id,
			"username": account.Username,
			"email":    account.Email,
			// Không trả về password
		},
	})
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	accountCollection := config.DB.Collection("accounts")
	if accountCollection == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database collection not initialized")
		return
	}

	var loginData models.Account
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Log để debug
	fmt.Printf("Login attempt for email: %s\n", loginData.Email)

	var account models.Account
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tìm account theo email và password
	err = accountCollection.FindOne(ctx, bson.M{
		"email":    loginData.Email,
		"password": loginData.Password,
	}).Decode(&account)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log để debug
	fmt.Printf("Login successful for user: %s\n", account.Username)

	// Trả về thông tin user
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user": map[string]interface{}{
			"id":       account.Id,
			"username": account.Username,
			"email":    account.Email,
		},
	})
}
