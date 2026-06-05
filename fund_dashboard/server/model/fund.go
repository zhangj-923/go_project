package model

import "time"

// Fund 基金基本信息
type Fund struct {
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // 例如：股票型、混合型等
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
