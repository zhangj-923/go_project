package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RealTimeNav struct {
	FundCode string `json:"fundcode"`
	Name     string `json:"name"`
	Jzrq     string `json:"jzrq"`   // 净值日期
	Dwjz     string `json:"dwjz"`   // 单位净值
	Gsz      string `json:"gsz"`    // 估算值
	Gszzl    string `json:"gszzl"`  // 估算增长率
	Gztime   string `json:"gztime"` // 估值时间
}

// GetRealTimeNav 获取实时基金估值
func GetRealTimeNav(code string) (*RealTimeNav, error) {
	url := fmt.Sprintf("http://fundgz.1234567.com.cn/js/%s.js?rt=%d", code, time.Now().UnixMilli())

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	strBody := string(body)
	if !strings.HasPrefix(strBody, "jsonpgz(") {
		return nil, fmt.Errorf("unexpected response format: %s", strBody)
	}

	jsonStr := strings.TrimSuffix(strings.TrimPrefix(strBody, "jsonpgz("), ");")
	var nav RealTimeNav
	if err := json.Unmarshal([]byte(jsonStr), &nav); err != nil {
		return nil, err
	}

	return &nav, nil
}

type NavHistoryResponse struct {
	Data struct {
		LSJZList []struct {
			FSRQ  string `json:"FSRQ"`  // 净值日期
			DWJZ  string `json:"DWJZ"`  // 单位净值
			LJJZ  string `json:"LJJZ"`  // 累计净值
			JZZZL string `json:"JZZZL"` // 净值增长率
		} `json:"LSJZList"`
	} `json:"Data"`
	ErrCode int `json:"ErrCode"`
}

// GetNavHistory 获取历史净值
func GetNavHistory(code string, pageIndex, pageSize int) (*NavHistoryResponse, error) {
	url := fmt.Sprintf("https://api.fund.eastmoney.com/f10/lsjz?fundCode=%s&pageIndex=%d&pageSize=%d", code, pageIndex, pageSize)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://fundf10.eastmoney.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data NavHistoryResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if data.ErrCode != 0 {
		return nil, fmt.Errorf("api error code: %d", data.ErrCode)
	}

	return &data, nil
}
