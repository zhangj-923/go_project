package exporter

import (
	"demo/mysql-agent/internal/models"
	"log"
)

// VMStorageExporter 实现了 Exporter 接口，用于将数据发送到 VictoriaMetrics.
type VMStorageExporter struct {
	// 在这里可以添加 VictoriaMetrics 的 URL 等配置
}

// NewVMStorageExporter 创建一个新的 VMStorageExporter.
func NewVMStorageExporter() *VMStorageExporter {
	return &VMStorageExporter{}
}

// Export 模拟将指标发送到 VictoriaMetrics.
func (e *VMStorageExporter) Export(metrics []models.Metric) error {
	// TODO: 在这里实现 VictoriaMetrics 的数据写入逻辑.
	// VictoriaMetrics 兼容 Prometheus Remote Write, 所以可以复用逻辑.
	log.Printf("[VMStorageExporter] 成功发送 %d 条指标数据", len(metrics))
	return nil
}
