package main

import (
	"log"
	"net/http"

	"flying-server/handler"
)

func main() {
	hub := handler.NewHub()
	go hub.Run()

	// WebSocket 路由
	http.HandleFunc("/ws", hub.HandleWebSocket)

	// 静态文件服务（前端）
	http.Handle("/", http.FileServer(http.Dir("../web/dist")))

	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
