package exporter

import "demo/mysql-agent/internal/models"

// Exporter 是数据发送器的通用接口.
// 任何实现了 Export() 方法的结构体都可以作为一个发送器.
type Exporter interface {
	Export([]models.Metric) error
}
