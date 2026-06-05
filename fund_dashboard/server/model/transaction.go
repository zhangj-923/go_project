package model

import "time"

// Transaction 交易记录
type Transaction struct {
	ID        int64     `json:"id"`
	FundCode  string    `json:"fund_code"`
	Type      string    `json:"type"` // "buy" 或 "sell"
	Date      string    `json:"date"` // YYYY-MM-DD
	Shares    float64   `json:"shares"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"` // 交易金额
	Fee       float64   `json:"fee"`    // 手续费
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
