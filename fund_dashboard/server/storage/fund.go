package storage

import (
	"fund_dashboard/model"
)

// GetAllFunds 获取所有基金
func GetAllFunds() ([]model.Fund, error) {
	rows, err := DB.Query("SELECT code, name, type, created_at, updated_at FROM funds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var funds []model.Fund
	for rows.Next() {
		var f model.Fund
		if err := rows.Scan(&f.Code, &f.Name, &f.Type, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		funds = append(funds, f)
	}

	// 如果为空返回空切片而不是nil
	if funds == nil {
		funds = make([]model.Fund, 0)
	}
	return funds, nil
}

// AddFund 添加基金
func AddFund(f *model.Fund) error {
	_, err := DB.Exec("INSERT INTO funds (code, name, type) VALUES (?, ?, ?)", f.Code, f.Name, f.Type)
	return err
}

// DeleteFund 删除基金
func DeleteFund(code string) error {
	_, err := DB.Exec("DELETE FROM funds WHERE code = ?", code)
	return err
}
