// go-product-api/routes/product_routes.go
package routes

import (
	"My_WEB/controllers"

	"github.com/gorilla/mux"
)

// ProductRoutes đăng ký các endpoint API cho Product
func ProductRoutes(router *mux.Router) {
	// Endpoint để tạo sản phẩm mới (POST /products)
	router.HandleFunc("/products", controllers.CreateProduct).Methods("POST")
	// Endpoint để lấy tất cả sản phẩm (GET /products)
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	// Endpoint để lấy sản phẩm theo ID (GET /products/{id})
	router.HandleFunc("/products/{id}", controllers.GetProductByID).Methods("GET")
	// Endpoint để cập nhật sản phẩm theo ID (PUT /products/{id})
	router.HandleFunc("/products/{id}", controllers.UpdateProduct).Methods("PUT")
	// Endpoint để xóa sản phẩm theo ID (DELETE /products/{id})
	router.HandleFunc("/products/{id}", controllers.DeleteProduct).Methods("DELETE")
}
func UserRoutes(router *mux.Router) {
	router.HandleFunc("/signup", controllers.SignUp).Methods("POST")
	router.HandleFunc("/signin", controllers.SignIn).Methods("POST")

}
