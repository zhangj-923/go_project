package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL

	//_ "github.com/lib/pq" // PostgreSQL 驱动
	"github.com/vito-go/sqlview"
)

func main() {
	// 从环境变量中读取数据库连接字符串，避免硬编码敏感信息
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("环境变量 DB_DSN 未设置")
	}

	// 连接到数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close() // 确保在 main 函数退出时关闭数据库连接

	// 设置数据库连接池参数（可选，但推荐）
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// 创建 SQLView 实例（自动检测数据库类型）
	viewer := sqlview.New(db, "/sqlview")

	// 挂载到 HTTP 服务器
	mux := http.NewServeMux()
	viewer.Mount(mux)

	log.Println("SQLView 运行于 http://localhost:8080/sqlview")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
