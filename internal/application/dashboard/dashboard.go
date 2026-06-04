package dashboard

import (
	"context"
	"sync"

	"github.com/juanpabloaj/watchdogdashboard/internal/domain"
)

const todosUnavailableWarning = "Todos Unavailable"

type UserRepository interface {
	GetUser(ctx context.Context, id int) (*domain.User, error)
}

type TodoRepository interface {
	GetTodosByUserID(ctx context.Context, userID int) ([]*domain.Todo, error)
}

type Service struct {
	userRepo UserRepository
	todoRepo TodoRepository
}

func NewService(userRepo UserRepository, todoRepo TodoRepository) *Service {
	return &Service{
		userRepo: userRepo,
		todoRepo: todoRepo,
	}
}

func (s *Service) GetDashboard(ctx context.Context, userID int) (*domain.Dashboard, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var user *domain.User
	var todos []*domain.Todo
	var userErr error
	var todosErr error

	go func() {
		defer wg.Done()
		u, err := s.userRepo.GetUser(ctx, userID)
		user = u
		userErr = err
	}()

	go func() {
		defer wg.Done()
		t, err := s.todoRepo.GetTodosByUserID(ctx, userID)
		todos = t
		todosErr = err
	}()

	wg.Wait()

	if userErr != nil {
		return nil, userErr
	}

	agg := aggregate(user, todos)

	if todosErr != nil {
		agg.ErrorWarning = stringPtr(todosUnavailableWarning)
	}

	return agg, nil
}

func aggregate(user *domain.User, todos []*domain.Todo) *domain.Dashboard {
	status := "Rookie"
	if user.Age > 50 {
		status = "Veteran"
	}

	pendings := []*domain.Todo{}
	for _, todo := range todos {
		if todo.Completed {
			continue
		}

		pendings = append(pendings, todo)
	}

	nextUrgentTask := ""
	if len(pendings) > 0 {
		nextUrgentTask = pendings[0].Todo
	}

	return &domain.Dashboard{
		ID:               user.ID,
		FullName:         user.FirstName + " " + user.LastName,
		Status:           status,
		PendingTaskCount: len(pendings),
		NextUrgentTask:   nextUrgentTask,
	}
}

func stringPtr(s string) *string {
	return &s
}
