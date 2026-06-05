package collector

import (
	"database/sql"
	"demo/mysql-agent/internal/models"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLCollector 实现了 Collector 接口，用于采集 MySQL 性能数据.
type MySQLCollector struct {
	db *sql.DB
}

// NewMySQLCollector 创建并返回一个 MySQLCollector 实例.
func NewMySQLCollector(dsn string) (*MySQLCollector, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("无法连接到 MySQL: %w", err)
	}
	return &MySQLCollector{db: db}, nil
}

// Collect 从 MySQL 的 'SHOW GLOBAL STATUS' 中采集指标.
func (c *MySQLCollector) Collect() ([]models.Metric, error) {
	rows, err := c.db.Query("SHOW GLOBAL STATUS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.Metric
	for rows.Next() {
		var name string
		var valueStr string
		if err := rows.Scan(&name, &valueStr); err != nil {
			continue
		}

		var val float64
		n, err := fmt.Sscanf(valueStr, "%f", &val)
		if err == nil && n == 1 {
			metrics = append(metrics, models.Metric{
				Name:      fmt.Sprintf("mysql_global_status_%s", name),
				Value:     val,
				Labels:    map[string]string{"component": "mysql"},
				Timestamp: time.Now(),
			})
		}
	}
	return metrics, nil
}
