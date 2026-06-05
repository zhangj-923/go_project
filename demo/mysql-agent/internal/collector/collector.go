package collector

import "demo/mysql-agent/internal/models"

// Collector 是数据采集器的通用接口.
// 任何实现了 Collect() 方法的结构体都可以作为一个采集器.
type Collector interface {
	Collect() ([]models.Metric, error)
}
