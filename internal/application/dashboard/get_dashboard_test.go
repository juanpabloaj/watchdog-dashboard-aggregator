package dashboard

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/juanpabloaj/watchdogdashboard/internal/domain"
)

type mockUserRepo struct {
	User *domain.User
	Err  error
}

func (m *mockUserRepo) GetUser(ctx context.Context, id int) (*domain.User, error) {
	return m.User, m.Err
}

type mockTodoRepo struct {
	Todos []*domain.Todo
	Err   error
}

func (m *mockTodoRepo) GetTodosByUserID(ctx context.Context, userID int) ([]*domain.Todo, error) {
	return m.Todos, m.Err
}

func TestGetDashboard(t *testing.T) {
	mockUser := &mockUserRepo{
		User: &domain.User{
			ID:        1,
			FirstName: "Terry",
			LastName:  "Medhurst",
			Age:       50,
		},
		Err: nil,
	}
	mockTodo := &mockTodoRepo{
		Todos: []*domain.Todo{
			{
				ID:        1,
				Todo:      "Do something important",
				Completed: false,
			},
			{
				ID:        2,
				Todo:      "Do something important later",
				Completed: false,
			},
		},
		Err: nil,
	}

	expected := &domain.Dashboard{
		ID:               1,
		FullName:         "Terry Medhurst",
		Status:           "Rookie",
		PendingTaskCount: 2,
		NextUrgentTask:   "Do something important",
		ErrorWarning:     nil,
	}

	tt := []struct {
		name     string
		expected *domain.Dashboard
		user     *mockUserRepo
		todo     *mockTodoRepo
	}{
		{"happy path", expected, mockUser, mockTodo},
		{
			"user ok todos error, should return warning",
			&domain.Dashboard{
				ID:               1,
				FullName:         "Terry Medhurst",
				Status:           "Rookie",
				PendingTaskCount: 0,
				NextUrgentTask:   "",
				ErrorWarning:     stringPtr(todosUnavailableWarning),
			},
			mockUser,
			&mockTodoRepo{
				Err: errors.New("mock todos error"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := NewService(tc.user, tc.todo)

			result, err := service.GetDashboard(context.Background(), 1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetDashboardUserRepoError(t *testing.T) {
	mockTodo := &mockTodoRepo{
		Todos: []*domain.Todo{
			{
				ID:        1,
				Todo:      "Do something important",
				Completed: false,
			},
			{
				ID:        2,
				Todo:      "Do something important later",
				Completed: false,
			},
		},
		Err: nil,
	}

	mockUserErr := errors.New("mock user error")
	mockUser := &mockUserRepo{
		Err: mockUserErr,
	}

	service := NewService(mockUser, mockTodo)

	result, err := service.GetDashboard(context.Background(), 1)
	if !errors.Is(err, mockUserErr) {
		t.Fatalf("unexpected error: got %v, want %v", err, mockUserErr)
	}

	if result != nil {
		t.Errorf("expected result to be nil, got %v", result)
	}
}
