package handler

import (
	"fund_dashboard/service"
	"net/http"
)

func GetPortfolio(w http.ResponseWriter, r *http.Request) {
	summary, err := service.CalculatePortfolio()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}
