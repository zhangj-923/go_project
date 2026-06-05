package main

import (
	"demo/mysql-agent/internal/agent"
	"demo/mysql-agent/internal/collector"
	"demo/mysql-agent/internal/exporter"
	"demo/mysql-agent/internal/processor"
	"log"
	"time"
)

func main() {
	// ====================================
	// 1. 配置与初始化
	// ====================================

	// MySQL 数据源配置
	dsn := "root:123456@tcp(127.0.0.1:3306)/"

	// Agent 采集间隔
	scrapeInterval := 15 * time.Second

	// ====================================
	// 2. 依赖注入与组装
	// ====================================

	// 初始化采集器
	mysqlCollector, err := collector.NewMySQLCollector(dsn)
	if err != nil {
		log.Fatalf("初始化 MySQL 采集器失败: %v", err)
	}

	// 初始化处理器
	govProcessor := processor.NewDataGovernanceProcessor()

	// 初始化发送器
	promExporter := exporter.NewPrometheusExporter()
	vmExporter := exporter.NewVMStorageExporter()

	// 创建 Agent 实例
	app := agent.NewAgent(scrapeInterval)

	// 将组件注册到 Agent
	app.AddCollector(mysqlCollector)
	app.AddProcessor(govProcessor)
	app.AddExporter(promExporter)
	app.AddExporter(vmExporter)

	// ====================================
	// 3. 启动 Agent
	// ====================================
	log.Println("MySQL Agent 启动中...")
	app.Start()
}
