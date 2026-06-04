package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/juanpabloaj/watchdogdashboard/internal/domain"
)

type mockDashboardService struct {
	Dashboard *domain.Dashboard
	Error     error
}

func (m *mockDashboardService) GetDashboard(ctx context.Context, userID int) (*domain.Dashboard, error) {
	return m.Dashboard, m.Error
}

func TestGetDashboard(t *testing.T) {
	expected := &domain.Dashboard{
		ID:               1,
		FullName:         "Terry Medhurst",
		Status:           "Rookie",
		PendingTaskCount: 2,
		NextUrgentTask:   "Do something important",
		ErrorWarning:     nil,
	}

	mock := &mockDashboardService{
		Dashboard: expected,
	}

	dhandler := NewDashboardHandler(mock)

	req, err := http.NewRequest(http.MethodGet, "/dashboard/1", nil)
	if err != nil {
		t.Fatalf("error creating request %v", err)
	}

	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(dhandler.GetDashboard)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var result domain.Dashboard
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("error decoding response %v", err)
	}

	if diff := cmp.Diff(expected, &result); diff != "" {
		t.Errorf("unexpected body: got %v want %v\ndiff %s", result, expected, diff)
	}
}

func TestGetDashboardError(t *testing.T) {
	errMockService := errors.New("mock service error")
	mock := &mockDashboardService{
		Error: errMockService,
	}

	dhandler := NewDashboardHandler(mock)

	req, err := http.NewRequest(http.MethodGet, "/dashboard/2", nil)
	if err != nil {
		t.Fatalf("error creating request %v", err)
	}

	req.SetPathValue("id", "2")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(dhandler.GetDashboard)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadGateway {
		t.Errorf("wrong status code: got %v want %v",
			status, http.StatusBadGateway)
	}

	if rr.Body.String() != "failed to fetch dashboard data\n" {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), "failed to fetch dashboard data\n")
	}
}
