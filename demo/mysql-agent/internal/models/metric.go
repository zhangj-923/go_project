package models

import "time"

// Metric 统一指标数据模型
// 这是整个 agent 数据流中流转的核心结构体.
type Metric struct {
	Name      string
	Value     float64
	Labels    map[string]string
	Timestamp time.Time
}
