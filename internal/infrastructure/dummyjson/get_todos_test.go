package dummyjson

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/juanpabloaj/watchdogdashboard/internal/domain"
)

func TestGetTodosByUserID(t *testing.T) {
	response := `{
"todos": [
{
"id": 47,
"todo": "Learn Javascript",
"completed": false,
"userId": 1
},
{
"id": 64,
"todo": "Listen to a new music genre",
"completed": true,
"userId": 1
}
],
"total": 2,
"skip": 0,
"limit": 2
}`

	expected := []*domain.Todo{
		{
			ID:        47,
			Todo:      "Learn Javascript",
			Completed: false,
		},
		{
			ID:        64,
			Todo:      "Listen to a new music genre",
			Completed: true,
		},
	}

	tt := []struct {
		name     string
		response string
		expected []*domain.Todo
	}{
		{
			"happy path",
			response,
			expected,
		},
		{
			"empty todo list",
			`{
"todos": [],
"total": 2,
"skip": 0,
"limit": 2
}`,
			[]*domain.Todo{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/todos/user/1" {
					t.Fatalf("unexpected URL path: %s", r.URL.Path)
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tc.response))
			}))
			defer server.Close()

			client := NewDummyJSONClient(server.URL, server.Client())

			resp, err := client.GetTodosByUserID(context.Background(), 1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetTodosByUserIDUnexpectedStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/todos/user/1" {
			t.Fatalf("unexpected URL path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(``))
	}))
	defer server.Close()

	client := NewDummyJSONClient(server.URL, server.Client())

	resp, err := client.GetTodosByUserID(context.Background(), 1)
	if err == nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !errors.Is(err, ErrUnexpectedStatusCode) {
		t.Errorf("expected status code error: %v", err)
	}

	if resp != nil {
		t.Errorf("expected nil response, got: %v", resp)
	}
}
