package handler

import (
	"encoding/json"
	"net/http"

	"fund_dashboard/service"
)

// ImportScreenshotHandler 处理截图上传
func ImportScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	// 仅支持 POST 方法
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 限制上传文件大小为 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse multipart form"})
		return
	}

	// 获取文件
	file, _, err := r.FormFile("file") // 假设前端字段名为 "file"
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// 解析截图
	funds, err := service.ParseScreenshot(file)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(funds)
}
