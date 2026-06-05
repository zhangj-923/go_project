package handler

import (
	"encoding/json"
	"fund_dashboard/model"
	"fund_dashboard/storage"
	"net/http"
)

func GetFunds(w http.ResponseWriter, r *http.Request) {
	funds, err := storage.GetAllFunds()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, funds)
}

func AddFund(w http.ResponseWriter, r *http.Request) {
	var f model.Fund
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := storage.AddFund(&f); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, f)
}

func DeleteFund(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "code is required"})
		return
	}

	if err := storage.DeleteFund(code); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
