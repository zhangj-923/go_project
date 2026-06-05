package exporter

import (
	"demo/mysql-agent/internal/models"
	"log"
)

// PrometheusExporter 实现了 Exporter 接口，用于将数据发送到 Prometheus.
type PrometheusExporter struct {
	// 在这里可以添加 Prometheus Remote Write 的 URL 等配置
}

// NewPrometheusExporter 创建一个新的 PrometheusExporter.
func NewPrometheusExporter() *PrometheusExporter {
	return &PrometheusExporter{}
}

// Export 模拟将指标发送到 Prometheus.
func (e *PrometheusExporter) Export(metrics []models.Metric) error {
	// TODO: 在这里实现 Prometheus Remote Write 协议的逻辑.
	// 为了演示，我们只打印日志.
	log.Printf("[PrometheusExporter] 成功发送 %d 条指标数据", len(metrics))
	return nil
}
