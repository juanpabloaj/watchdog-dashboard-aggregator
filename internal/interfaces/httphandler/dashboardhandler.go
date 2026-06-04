package httphandler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/juanpabloaj/watchdogdashboard/internal/domain"
)

const (
	minUserID = 1
	maxUserID = 8
)

type DashboardService interface {
	GetDashboard(ctx context.Context, userID int) (*domain.Dashboard, error)
}

type DashboardHandler struct {
	service DashboardService
}

func NewDashboardHandler(service DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if userID < minUserID || userID > maxUserID {
		http.Error(w, "user id must be between 1 and 8", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetDashboard(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch dashboard data", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
