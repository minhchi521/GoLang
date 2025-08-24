// go-product-api/utils/response.go
package utils

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON gửi phản hồi JSON với mã trạng thái HTTP cụ thể
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)               // Chuyển đổi payload (dữ liệu Go) sang JSON byte array
	w.Header().Set("Content-Type", "application/json") // Đặt Content-Type header
	w.WriteHeader(code)                                // Đặt mã trạng thái HTTP
	w.Write(response)                                  // Ghi dữ liệu JSON vào phản hồi
}

// RespondWithError gửi phản hồi lỗi JSON với mã trạng thái HTTP cụ thể
func RespondWithError(w http.ResponseWriter, code int, message string) {
	// Tạo một map để chứa thông báo lỗi và gửi đi
	RespondWithJSON(w, code, map[string]string{"error": message})
}
