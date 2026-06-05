package processor

import "demo/mysql-agent/internal/models"

// Processor 是数据治理/处理器的通用接口.
// 任何实现了 Process() 方法的结构体都可以作为一个处理器.
type Processor interface {
	Process([]models.Metric) []models.Metric
}
