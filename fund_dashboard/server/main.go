package main

import (
	"fund_dashboard/handler"
	"fund_dashboard/storage"
	"log"
	"net/http"
)

func main() {
	// 初始化数据库
	storage.InitDB("data/fund.db")

	mux := http.NewServeMux()

	// 基金管理 API
	mux.HandleFunc("/api/funds", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetFunds(w, r)
		case http.MethodPost:
			handler.AddFund(w, r)
		case http.MethodDelete:
			handler.DeleteFund(w, r)
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// 交易记录 API
	mux.HandleFunc("/api/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTransactions(w, r)
		case http.MethodPost:
			handler.AddTransaction(w, r)
		case http.MethodDelete:
			handler.DeleteTransaction(w, r)
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// 投资组合分析 API
	mux.HandleFunc("/api/analytics/summary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetPortfolio(w, r)
		} else if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// 截图导入 API
	mux.HandleFunc("/api/import/screenshot", handler.ImportScreenshotHandler)

	// 添加 CORS 中间件
	handlerWithCORS := handler.EnableCORS(mux.ServeHTTP)

	log.Println("服务器运行在 :8080 端口")
	if err := http.ListenAndServe(":8080", handlerWithCORS); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
