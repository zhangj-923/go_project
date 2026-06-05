package service

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type ParsedFund struct {
	Code            string  `json:"code"`
	CurrentValue    float64 `json:"current_value"`
	TotalPnl        float64 `json:"total_pnl"`
	EstimatedCost   float64 `json:"estimated_cost"`
	EstimatedShares float64 `json:"estimated_shares"`
	Nav             float64 `json:"nav"`
}

// ParseScreenshot 解析截图并返回基金数据
func ParseScreenshot(fileData io.Reader) ([]ParsedFund, error) {
	apiKey := os.Getenv("VISION_API_KEY")
	if apiKey == "" {
		// 环境变量未设置，返回mock数据
		return getMockParsedFunds()
	}

	// 实际调用Vision API的逻辑可在此实现
	return nil, fmt.Errorf("vision API integration not implemented yet")
}

func getMockParsedFunds() ([]ParsedFund, error) {
	mocks := []struct {
		Code         string
		CurrentValue float64
		TotalPnl     float64
	}{
		{"014868", 6860.63, 446.34},
		{"008889", 5888.60, 1583.57},
		{"018447", 4030.36, -969.64},
		{"018957", 2439.41, 57.28},
	}

	var results []ParsedFund
	for _, m := range mocks {
		pf := ParsedFund{
			Code:         m.Code,
			CurrentValue: m.CurrentValue,
			TotalPnl:     m.TotalPnl,
		}

		// 计算估算成本
		pf.EstimatedCost = pf.CurrentValue - pf.TotalPnl

		// 获取实时净值或单位净值
		navData, err := GetRealTimeNav(pf.Code)
		if err == nil && navData != nil {
			var nav float64
			// 优先使用估算值
			if navData.Gsz != "" {
				nav, _ = strconv.ParseFloat(navData.Gsz, 64)
			}
			// 如果估算值无效，退而求其次使用单位净值
			if nav <= 0 && navData.Dwjz != "" {
				nav, _ = strconv.ParseFloat(navData.Dwjz, 64)
			}

			if nav > 0 {
				pf.Nav = nav
				pf.EstimatedShares = pf.CurrentValue / nav
			} else {
				// 兜底逻辑
				pf.Nav = 1.0
				pf.EstimatedShares = pf.CurrentValue
			}
		} else {
			// 如果调用API失败，使用兜底逻辑
			pf.Nav = 1.0
			pf.EstimatedShares = pf.CurrentValue
		}

		results = append(results, pf)
	}

	return results, nil
}
