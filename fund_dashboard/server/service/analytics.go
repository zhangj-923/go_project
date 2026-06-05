package service

import (
	"fund_dashboard/model"
	"fund_dashboard/storage"
	"strconv"
)

// CalculatePortfolio 计算投资组合汇总和个基持仓
func CalculatePortfolio() (*model.PortfolioSummary, error) {
	funds, err := storage.GetAllFunds()
	if err != nil {
		return nil, err
	}

	summary := &model.PortfolioSummary{
		Positions: make([]model.FundPosition, 0),
	}

	for _, f := range funds {
		txs, err := storage.GetTransactionsByFund(f.Code)
		if err != nil {
			return nil, err
		}

		if len(txs) == 0 {
			continue
		}

		var totalShares float64
		var totalCost float64

		for _, t := range txs {
			if t.Type == "buy" {
				totalShares += t.Shares
				totalCost += t.Amount
			} else if t.Type == "sell" {
				// 对于卖出，减少份额，同时按比例减少成本
				if totalShares > 0 {
					ratio := t.Shares / totalShares
					totalShares -= t.Shares
					totalCost -= totalCost * ratio
				}
			}
		}

		if totalShares <= 0.001 { // 持仓为0，忽略
			continue
		}

		avgCost := totalCost / totalShares

		// 获取实时净值或最新净值
		navData, err := GetRealTimeNav(f.Code)
		var currentNav float64
		var dailyPnLRate float64
		var navDate string
		var estNavDate string

		if err == nil && navData != nil {
			currentNav, _ = strconv.ParseFloat(navData.Gsz, 64)
			dailyPnLRate, _ = strconv.ParseFloat(navData.Gszzl, 64)
			navDate = navData.Jzrq
			estNavDate = navData.Gztime
		} else {
			// 退化到获取历史净值
			histNav, err := GetNavHistory(f.Code, 1, 1)
			if err == nil && histNav != nil && len(histNav.Data.LSJZList) > 0 {
				currentNav, _ = strconv.ParseFloat(histNav.Data.LSJZList[0].DWJZ, 64)
				dailyPnLRate, _ = strconv.ParseFloat(histNav.Data.LSJZList[0].JZZZL, 64)
				navDate = histNav.Data.LSJZList[0].FSRQ
			}
		}

		currentValue := totalShares * currentNav
		// 今日盈亏粗略估算：当前总价值 * 今日涨跌幅 / (1 + 今日涨跌幅)
		// 或者用上个交易日价值来算，简化处理
		yesterdayValue := currentValue / (1 + dailyPnLRate/100.0)
		dailyPnL := currentValue - yesterdayValue

		totalPnL := currentValue - totalCost
		totalPnLRate := 0.0
		if totalCost > 0 {
			totalPnLRate = (totalPnL / totalCost) * 100.0
		}

		pos := model.FundPosition{
			FundCode:     f.Code,
			FundName:     f.Name,
			TotalShares:  totalShares,
			TotalCost:    totalCost,
			AverageCost:  avgCost,
			CurrentNav:   currentNav,
			CurrentValue: currentValue,
			DailyPnL:     dailyPnL,
			DailyPnLRate: dailyPnLRate,
			TotalPnL:     totalPnL,
			TotalPnLRate: totalPnLRate,
			NavDate:      navDate,
			EstNavDate:   estNavDate,
		}

		summary.Positions = append(summary.Positions, pos)
		summary.TotalCost += totalCost
		summary.TotalValue += currentValue
		summary.DailyPnL += dailyPnL
	}

	summary.TotalPnL = summary.TotalValue - summary.TotalCost
	if summary.TotalCost > 0 {
		summary.TotalPnLRate = (summary.TotalPnL / summary.TotalCost) * 100.0
	}
	if summary.TotalValue-summary.DailyPnL > 0 {
		summary.DailyPnLRate = (summary.DailyPnL / (summary.TotalValue - summary.DailyPnL)) * 100.0
	}

	return summary, nil
}
