package dashboard

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/juanpabloaj/watchdogdashboard/internal/domain"
)

func TestAggregate(t *testing.T) {
	expected := &domain.Dashboard{
		ID:               1,
		FullName:         "Terry Medhurst",
		Status:           "Rookie",
		PendingTaskCount: 2,
		NextUrgentTask:   "Do something important",
		ErrorWarning:     nil,
	}

	user := &domain.User{
		ID:        1,
		FirstName: "Terry",
		LastName:  "Medhurst",
		Age:       50,
	}

	todos := []*domain.Todo{
		{
			Id:        1,
			Todo:      "Do something important",
			Completed: false,
		},
		{
			Id:        2,
			Todo:      "Do something important later",
			Completed: false,
		},
		{
			Completed: true,
		},
	}

	tt := []struct {
		name     string
		expected *domain.Dashboard
		user     *domain.User
		todos    []*domain.Todo
	}{
		{
			"happy path",
			expected,
			user,
			todos,
		},
		{
			"veteran status",
			&domain.Dashboard{
				ID:               1,
				FullName:         "Terry Medhurst",
				Status:           "Veteran",
				PendingTaskCount: 2,
				NextUrgentTask:   "Do something important",
				ErrorWarning:     nil,
			},
			&domain.User{
				ID:        1,
				FirstName: "Terry",
				LastName:  "Medhurst",
				Age:       51,
			},
			todos,
		},
		{
			"empty todos",
			&domain.Dashboard{
				ID:               1,
				FullName:         "Terry Medhurst",
				Status:           "Veteran",
				PendingTaskCount: 0,
				NextUrgentTask:   "",
				ErrorWarning:     nil,
			},
			&domain.User{
				ID:        1,
				FirstName: "Terry",
				LastName:  "Medhurst",
				Age:       51,
			},
			[]*domain.Todo{},
		},
		// TODO Todos Unavailable warning
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := aggregate(tc.user, tc.todos)
			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
