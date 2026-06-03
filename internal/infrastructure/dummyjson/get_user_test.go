package dummyjson

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/juanpabloaj/watchdogdashboard/internal/domain"
)

func TestGetUser(t *testing.T) {
	response := `
	{
"id": 1,
"firstName": "Emily",
"lastName": "Johnson",
"age": 29}`

	expected := &domain.User{
		ID:        1,
		FirstName: "Emily",
		LastName:  "Johnson",
		Age:       29,
	}

	tt := []struct {
		name     string
		response string
		expected *domain.User
	}{
		{
			"happy path",
			response,
			expected,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/users/1" {
					t.Fatalf("unexpected URL path: %s", r.URL.Path)
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tc.response))
			}))
			defer server.Close()

			client := NewDummyJSONClient(server.URL, server.Client())

			resp, err := client.GetUser(context.Background(), 1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetUserUnexpectedStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/1" {
			t.Fatalf("unexpected URL path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(``))
	}))
	defer server.Close()

	client := NewDummyJSONClient(server.URL, server.Client())

	resp, err := client.GetUser(context.Background(), 1)
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

func TestGetUserTimeoutWithSleep(t *testing.T) {
	t.Skip("skipping it because is slow, to try timeout with sleep comment this line")
	response := `
	{
"id": 1,
"firstName": "Emily",
"lastName": "Johnson",
"age": 29}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/1" {
			t.Fatalf("unexpected URL path: %s", r.URL.Path)
		}

		// forcing timeout with sleep
		time.Sleep(3 * time.Second)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	// http client for test with 2 seconds timeout
	httpClient := server.Client()
	httpClient.Timeout = 2 * time.Second

	client := NewDummyJSONClient(server.URL, httpClient)

	resp, err := client.GetUser(context.Background(), 1)
	if err == nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected timeout error: %v", err)
	}

	if resp != nil {
		t.Errorf("expected nil response, got: %v", resp)
	}
}

type mockHTTPClientTimeout struct {
}

func (c *mockHTTPClientTimeout) Do(req *http.Request) (*http.Response, error) {
	return nil, context.DeadlineExceeded
}

func TestGetUserTimeoutWithMock(t *testing.T) {
	mockClient := &mockHTTPClientTimeout{}

	client := NewDummyJSONClient("dummyurl", mockClient)

	resp, err := client.GetUser(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected timeout error: %v", err)
	}

	if resp != nil {
		t.Errorf("expected nil response, got: %v", resp)
	}
}
