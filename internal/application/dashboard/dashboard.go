package dashboard

import (
	"context"

	"github.com/juanpabloaj/watchdogdashboard/internal/domain"
)

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

func NewDashboard(UserRepository, TodoRepository) *Service {
	return &Service{}
}

func (s *Service) GetDashboard(ctx context.Context, userID int) (*domain.Dashboard, error) {
	user := &domain.User{}
	todos := []*domain.Todo{}
	return aggregate(user, todos), nil
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
