package model

// FundPosition 单个基金持仓情况
type FundPosition struct {
	FundCode     string  `json:"fund_code"`
	FundName     string  `json:"fund_name"`
	TotalShares  float64 `json:"total_shares"`
	TotalCost    float64 `json:"total_cost"`
	AverageCost  float64 `json:"average_cost"`
	CurrentNav   float64 `json:"current_nav"`    // 当前净值或实时估值
	CurrentValue float64 `json:"current_value"`  // 当前市值
	DailyPnL     float64 `json:"daily_pnl"`      // 今日盈亏
	DailyPnLRate float64 `json:"daily_pnl_rate"` // 今日盈亏率
	TotalPnL     float64 `json:"total_pnl"`      // 累计盈亏
	TotalPnLRate float64 `json:"total_pnl_rate"` // 累计盈亏率
	NavDate      string  `json:"nav_date"`       // 净值日期
	EstNavDate   string  `json:"est_nav_date"`   // 估算时间
}

// PortfolioSummary 投资组合汇总情况
type PortfolioSummary struct {
	TotalCost    float64        `json:"total_cost"`
	TotalValue   float64        `json:"total_value"`
	DailyPnL     float64        `json:"daily_pnl"`
	DailyPnLRate float64        `json:"daily_pnl_rate"`
	TotalPnL     float64        `json:"total_pnl"`
	TotalPnLRate float64        `json:"total_pnl_rate"`
	Positions    []FundPosition `json:"positions"`
}
