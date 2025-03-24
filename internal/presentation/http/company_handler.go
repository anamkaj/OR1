package http

import (
	"encoding/json"
	"net/http"
	"ord_crm/internal/application"
)

// CompanyHandler - обработчик HTTP-запросов
type CompanyHandler struct {
	service *application.CompanyService
}

// NewCompanyHandler создает новый обработчик
func NewCompanyHandler(s *application.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: s}
}

// GetLotHandler - обработчик для получения компаний
func (h *CompanyHandler) GetLotHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id параметр отсутствует", http.StatusBadRequest)
		return
	}

	companies, err := h.service.GetClientList("", 2)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}
