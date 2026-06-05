package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Metric struct {
	Name      string
	Value     float64
	Labels    map[string]string
	Timestamp time.Time
}

// Collector 数据采集器接口
type Collector interface {
	Collect() ([]Metric, error)
}

// Processor 数据治理/处理器接口
type Processor interface {
	Process([]Metric) []Metric
}

// Exporter 数据发送器接口
type Exporter interface {
	Export([]Metric) error
}

// ==================
// 1. 采集器实现
// ==================

// MySQLCollector MySQL 性能数据采集器
type MySQLCollector struct {
	db *sql.DB
}

func NewMySQLCollector(dsn string) (*MySQLCollector, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &MySQLCollector{db: db}, nil
}

func (c *MySQLCollector) Collect() ([]Metric, error) {
	// 获取全局状态指标作为示例
	rows, err := c.db.Query("SHOW GLOBAL STATUS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var name string
		var valueStr string
		if err := rows.Scan(&name, &valueStr); err != nil {
			continue
		}

		// 尝试解析为数值型指标
		var val float64
		n, err := fmt.Sscanf(valueStr, "%f", &val)
		if err == nil && n == 1 {
			metrics = append(metrics, Metric{
				Name:      fmt.Sprintf("mysql_global_status_%s", name),
				Value:     val,
				Labels:    map[string]string{"component": "mysql"},
				Timestamp: time.Now(),
			})
		}
	}
	return metrics, nil
}

// ==================
// 2. 数据治理实现
// ==================

// DataGovernanceProcessor 负责清洗、标准化指标数据
type DataGovernanceProcessor struct {
}

func (p *DataGovernanceProcessor) Process(metrics []Metric) []Metric {
	var processed []Metric
	for _, m := range metrics {
		// 规则1：过滤异常负值
		if m.Value < 0 {
			continue
		}

		// 规则2：标准化命名（如转为小写）
		m.Name = strings.ToLower(m.Name)

		// 规则3：补充通用标签
		if m.Labels == nil {
			m.Labels = make(map[string]string)
		}
		m.Labels["env"] = "production"

		processed = append(processed, m)
	}
	return processed
}

// ==================
// 3. 数据发送器实现
// ==================

// PrometheusExporter 适配 Prometheus 的导出器（可使用 Remote Write 或拉取模型）
type PrometheusExporter struct {
	// 可在此处配置 Remote Write URL
}

func (e *PrometheusExporter) Export(metrics []Metric) error {
	// TODO: 实现 Prometheus Remote Write 协议，或者写入本地供 /metrics 接口拉取
	log.Printf("[PrometheusExporter] 成功发送 %d 条指标数据", len(metrics))
	return nil
}

// VMStorageExporter 适配 VictoriaMetrics 的导出器
type VMStorageExporter struct {
	// 可配置 VM URL
}

func (e *VMStorageExporter) Export(metrics []Metric) error {
	// TODO: VictoriaMetrics 支持 Prometheus Remote Write 协议及自己的协议
	log.Printf("[VMStorageExporter] 成功发送 %d 条指标数据", len(metrics))
	return nil
}

// ==================
// 4. Agent 核心调度机制
// ==================

// Agent 协调器
type Agent struct {
	Collectors []Collector
	Processors []Processor
	Exporters  []Exporter
	Interval   time.Duration
}

func (a *Agent) Start() {
	ticker := time.NewTicker(a.Interval)
	defer ticker.Stop()

	for range ticker.C {
		var allMetrics []Metric

		// 1. 数据采集 (Collect)
		for _, c := range a.Collectors {
			metrics, err := c.Collect()
			if err != nil {
				log.Printf("采集失败: %v", err)
				continue
			}
			allMetrics = append(allMetrics, metrics...)
		}

		// 2. 数据治理 (Process)
		for _, p := range a.Processors {
			allMetrics = p.Process(allMetrics)
		}

		// 3. 多端分发 (Export)
		var wg sync.WaitGroup
		for _, e := range a.Exporters {
			wg.Add(1)
			go func(exp Exporter) {
				defer wg.Done()
				if err := exp.Export(allMetrics); err != nil {
					log.Printf("发送失败: %v", err)
				}
			}(e)
		}
		wg.Wait()
	}
}

func main() {
	// 配置 MySQL DSN
	dsn := "water_manage:Zrzk123@@tcp(172.168.10.192:13306)/"

	collector, err := NewMySQLCollector(dsn)
	if err != nil {
		log.Fatalf("初始化 MySQL 采集器失败: %v", err)
	}

	// 组装 Agent 管道
	agent := &Agent{
		Collectors: []Collector{collector},
		Processors: []Processor{&DataGovernanceProcessor{}},
		Exporters: []Exporter{
			&PrometheusExporter{},
			&VMStorageExporter{},
		},
		Interval: 15 * time.Second,
	}

	log.Println("Agent 启动中...")
	agent.Start()
}
