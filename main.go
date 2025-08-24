package main

import (
	"log"
	"net/http"

	"My_WEB/config"
	"My_WEB/middleware"
	"My_WEB/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Kết nối đến database
	config.ConnectDB()

	// Test kết nối
	if config.DB == nil {
		log.Fatal("Database connection failed")
	}

	// Khởi tạo router
	router := mux.NewRouter()

	// Định nghĩa các route
	//router.HandleFunc("/products", controllers.CreateProduct).Methods("POST")
	routes.ProductRoutes(router)
	routes.UserRoutes(router)

	// Áp dụng CORS middleware
	corsRouter := middleware.EnableCORS(router)
	// Chạy server
	log.Fatal(http.ListenAndServe(":8080", corsRouter))
}
