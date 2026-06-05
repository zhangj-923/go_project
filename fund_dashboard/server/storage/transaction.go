package storage

import (
	"fund_dashboard/model"
)

// GetTransactionsByFund 获取特定基金的交易记录
func GetTransactionsByFund(fundCode string) ([]model.Transaction, error) {
	rows, err := DB.Query("SELECT id, fund_code, type, date, shares, price, amount, fee, created_at, updated_at FROM transactions WHERE fund_code = ? ORDER BY date DESC", fundCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.FundCode, &t.Type, &t.Date, &t.Shares, &t.Price, &t.Amount, &t.Fee, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	if txs == nil {
		txs = make([]model.Transaction, 0)
	}
	return txs, nil
}

// GetAllTransactions 获取所有交易记录
func GetAllTransactions() ([]model.Transaction, error) {
	rows, err := DB.Query("SELECT id, fund_code, type, date, shares, price, amount, fee, created_at, updated_at FROM transactions ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.FundCode, &t.Type, &t.Date, &t.Shares, &t.Price, &t.Amount, &t.Fee, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	if txs == nil {
		txs = make([]model.Transaction, 0)
	}
	return txs, nil
}

// AddTransaction 添加交易记录
func AddTransaction(t *model.Transaction) error {
	_, err := DB.Exec(`INSERT INTO transactions (fund_code, type, date, shares, price, amount, fee) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		t.FundCode, t.Type, t.Date, t.Shares, t.Price, t.Amount, t.Fee)
	return err
}

// DeleteTransaction 删除交易记录
func DeleteTransaction(id int64) error {
	_, err := DB.Exec("DELETE FROM transactions WHERE id = ?", id)
	return err
}
