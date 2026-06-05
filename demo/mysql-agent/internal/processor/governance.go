package processor

import (
	"demo/mysql-agent/internal/models"
	"strings"
)

// DataGovernanceProcessor 实现了 Processor 接口，负责清洗和标准化指标数据.
type DataGovernanceProcessor struct{}

// NewDataGovernanceProcessor 创建一个新的数据治理处理器.
func NewDataGovernanceProcessor() *DataGovernanceProcessor {
	return &DataGovernanceProcessor{}
}

// Process 应用一系列规则来处理指标.
func (p *DataGovernanceProcessor) Process(metrics []models.Metric) []models.Metric {
	var processed []models.Metric
	for _, m := range metrics {
		// 规则1：过滤异常负值
		if m.Value < 0 {
			continue
		}

		// 规则2：标准化命名（转为小写）
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
