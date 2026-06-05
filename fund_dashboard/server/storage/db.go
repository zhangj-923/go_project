package storage

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB 初始化SQLite数据库
func InitDB(dbPath string) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("创建数据库目录失败: %v", err)
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}

	createTables()
}

func createTables() {
	fundTable := `
	CREATE TABLE IF NOT EXISTS funds (
		code TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	transactionTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fund_code TEXT NOT NULL,
		type TEXT NOT NULL,
		date TEXT NOT NULL,
		shares REAL NOT NULL,
		price REAL NOT NULL,
		amount REAL NOT NULL,
		fee REAL NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (fund_code) REFERENCES funds(code)
	);`

	if _, err := DB.Exec(fundTable); err != nil {
		log.Fatalf("创建funds表失败: %v", err)
	}

	if _, err := DB.Exec(transactionTable); err != nil {
		log.Fatalf("创建transactions表失败: %v", err)
	}
}
