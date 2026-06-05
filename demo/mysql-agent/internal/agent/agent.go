package agent

import (
	"demo/mysql-agent/internal/collector"
	"demo/mysql-agent/internal/exporter"
	"demo/mysql-agent/internal/models"
	"demo/mysql-agent/internal/processor"
	"log"
	"sync"
	"time"
)

// Agent 是核心调度器，负责管理整个采集、处理、发送的生命周期.
type Agent struct {
	Collectors []collector.Collector
	Processors []processor.Processor
	Exporters  []exporter.Exporter
	Interval   time.Duration
}

// NewAgent 创建并配置一个新的 Agent.
func NewAgent(interval time.Duration) *Agent {
	return &Agent{
		Interval: interval,
	}
}

// AddCollector 添加一个采集器到 Agent.
func (a *Agent) AddCollector(c collector.Collector) {
	a.Collectors = append(a.Collectors, c)
}

// AddProcessor 添加一个处理器到 Agent.
func (a *Agent) AddProcessor(p processor.Processor) {
	a.Processors = append(a.Processors, p)
}

// AddExporter 添加一个发送器到 Agent.
func (a *Agent) AddExporter(e exporter.Exporter) {
	a.Exporters = append(a.Exporters, e)
}

// Start 启动 Agent 的主循环.
func (a *Agent) Start() {
	ticker := time.NewTicker(a.Interval)
	defer ticker.Stop()

	log.Println("Agent 核心调度器已启动...")

	for {
		select {
		case <-ticker.C:
			a.runOnce()
		}
	}
}

// runOnce 执行一次完整的采集、处理、发送流程.
func (a *Agent) runOnce() {
	var allMetrics []models.Metric
	var wg sync.WaitGroup

	// 1. 数据采集 (Collect)
	for _, c := range a.Collectors {
		wg.Add(1)
		go func(coll collector.Collector) {
			defer wg.Done()
			metrics, err := coll.Collect()
			if err != nil {
				log.Printf("采集器 %T 失败: %v", coll, err)
				return
			}
			// TODO: 需要考虑并发安全
			allMetrics = append(allMetrics, metrics...)
		}(c)
	}
	wg.Wait()

	// 2. 数据治理 (Process)
	for _, p := range a.Processors {
		allMetrics = p.Process(allMetrics)
	}

	// 3. 多端分发 (Export)
	for _, e := range a.Exporters {
		wg.Add(1)
		go func(exp exporter.Exporter) {
			defer wg.Done()
			if err := exp.Export(allMetrics); err != nil {
				log.Printf("发送器 %T 失败: %v", exp, err)
			}
		}(e)
	}
	wg.Wait()

	log.Printf("完成一轮处理, 共处理 %d 条指标", len(allMetrics))
}
