package handler

import (
	"encoding/json"
	"fund_dashboard/model"
	"fund_dashboard/storage"
	"net/http"
	"strconv"
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	fundCode := r.URL.Query().Get("fund_code")
	var txs []model.Transaction
	var err error

	if fundCode != "" {
		txs, err = storage.GetTransactionsByFund(fundCode)
	} else {
		txs, err = storage.GetAllTransactions()
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, txs)
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	var t model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := storage.AddTransaction(&t); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := storage.DeleteTransaction(id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
