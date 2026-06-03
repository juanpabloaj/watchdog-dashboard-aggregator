package dummyjson

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juanpabloaj/watchdogdashboard/internal/domain"
)

var ErrUnexpectedStatusCode = fmt.Errorf("unexpected status code")

type userResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

/*
{
"todos": [
{
"id": 47,
"todo": "Learn Javascript",
"completed": false,
"userId": 1
},
],
"total": 2,
"skip": 0,
"limit": 2
}
*/

type todosResponseData struct {
	Todos []Todo `json:"todos"`
}

type Todo struct {
	ID        int    `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	url        string
	httpClient HTTPClient
}

func NewDummyJSONClient(url string, httpClient HTTPClient) *Client {
	return &Client{
		url:        url,
		httpClient: httpClient,
	}
}

func (c *Client) GetUser(ctx context.Context, id int) (*domain.User, error) {
	url := fmt.Sprintf("%s/users/%d", c.url, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var user userResponse
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
	}, nil
}

func (c *Client) GetTodosByUserID(ctx context.Context, userID int) ([]*domain.Todo, error) {
	url := fmt.Sprintf("%s/todos/user/%d", c.url, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var todosResponse todosResponseData
	err = json.NewDecoder(resp.Body).Decode(&todosResponse)
	if err != nil {
		return nil, err
	}

	todos := make([]*domain.Todo, len(todosResponse.Todos))
	for i, todo := range todosResponse.Todos {
		todos[i] = &domain.Todo{
			ID:        todo.ID,
			Todo:      todo.Todo,
			Completed: todo.Completed,
		}
	}

	return todos, nil
}
