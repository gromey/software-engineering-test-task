package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"cruder/internal/controller"
	"cruder/internal/model"
	"cruder/internal/repository"
	"cruder/internal/service"
	"cruder/pkg/validation"

	"github.com/gin-gonic/gin"
)

type MockUserRepository struct {
	counter int64
	Users   []model.User
}

func (m *MockUserRepository) GetAll(_ context.Context) ([]model.User, error) {
	return m.Users, nil
}

func (m *MockUserRepository) GetByUsername(_ context.Context, username string) (*model.User, error) {
	for _, u := range m.Users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, validation.ErrUserNotFound
}

func (m *MockUserRepository) GetByID(_ context.Context, id int64) (*model.User, error) {
	for _, u := range m.Users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, validation.ErrUserNotFound
}

func (m *MockUserRepository) Post(_ context.Context, user *model.User) (int64, error) {
	m.counter++
	user.ID = m.counter
	m.Users = append(m.Users, *user)
	return user.ID, nil
}

func (m *MockUserRepository) Patch(_ context.Context, user *model.User) error {
	fmt.Println(user)
	for i, u := range m.Users {
		if u.ID == user.ID {
			m.Users[i] = *user
			return nil
		}
	}
	return validation.ErrUserNotFound
}

func (m *MockUserRepository) Delete(_ context.Context, id int64) error {
	for i, u := range m.Users {
		if u.ID == id {
			m.Users = append(m.Users[:i], m.Users[i+1:]...)
			break
		}
	}
	return nil
}

const testApiKey = "testApiKey"

var (
	user1 = model.User{ID: 1, Username: "jdoe", Email: "jdoe@example.com", FullName: "John Doe"}
	user2 = model.User{ID: 2, Username: "asmith", Email: "asmith@example.com", FullName: "Alice Smith"}
	user3 = model.User{ID: 3, Username: "bjones", Email: "bjones@example.com", FullName: "Bob Jones"}
)

func insertTestUser(mock *MockUserRepository, user *model.User) {
	_, _ = mock.Post(context.Background(), user)
}

func userExists(mock *MockUserRepository, id int64) bool {
	_, err := mock.GetByID(context.Background(), id)
	return err == nil
}

func requester(method, url string, body any, mockRepo repository.UserRepository) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	repositories := &repository.Repository{Users: mockRepo}
	services := service.NewService(repositories)
	controllers := controller.NewController(services)

	r := gin.Default()
	router := New(r, testApiKey, controllers.Users)

	var reader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reader = bytes.NewReader(data)
	}

	req, _ := http.NewRequest(method, url, reader)
	req.Header.Set("x-api-key", testApiKey)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func TestGetAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)

	insertTestUser(mockRepo, &user1)
	insertTestUser(mockRepo, &user2)
	insertTestUser(mockRepo, &user3)

	rr := requester(http.MethodGet, "/api/v1/users/", nil, mockRepo)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var users []model.User
	if err := json.Unmarshal(rr.Body.Bytes(), &users); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestGetUserByUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)

	insertTestUser(mockRepo, &user1)

	tests := []struct {
		name    string
		url     string
		expCode int
	}{
		{
			name:    "success",
			url:     "/api/v1/users/username/jdoe",
			expCode: http.StatusOK,
		},
		{
			name:    "user not found",
			url:     "/api/v1/users/username/asmith",
			expCode: http.StatusNotFound,
		},
		{
			name:    "bad request",
			url:     "/api/v1/users/username/a",
			expCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := requester(http.MethodGet, tt.url, nil, mockRepo)

			if rr.Code != tt.expCode {
				t.Errorf("expected status %d, got %d", tt.expCode, rr.Code)
			}

			if tt.expCode == http.StatusOK {
				var user model.User
				if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if !reflect.DeepEqual(user, user1) {
					t.Errorf("expected %v users, got %v", user1, user)
				}

			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)

	insertTestUser(mockRepo, &user1)

	tests := []struct {
		name    string
		url     string
		expCode int
	}{
		{
			name:    "success",
			url:     "/api/v1/users/id/1",
			expCode: http.StatusOK,
		},
		{
			name:    "user not found",
			url:     "/api/v1/users/id/2",
			expCode: http.StatusNotFound,
		},
		{
			name:    "bad request",
			url:     "/api/v1/users/id/0",
			expCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := requester(http.MethodGet, tt.url, nil, mockRepo)

			if rr.Code != tt.expCode {
				t.Errorf("expected status %d, got %d", tt.expCode, rr.Code)
			}

			if tt.expCode == http.StatusOK {
				var user model.User
				if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if !reflect.DeepEqual(user, user1) {
					t.Errorf("expected %v users, got %v", user1, user)
				}

			}
		})
	}
}

func TestPostUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	rr := requester(http.MethodPost, "/api/v1/users/", user1, mockRepo)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var id struct{ ID int64 }
	if err := json.Unmarshal(rr.Body.Bytes(), &id); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	rr = requester(http.MethodGet, fmt.Sprintf("/api/v1/users/id/%d", id.ID), nil, mockRepo)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var user model.User
	if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !reflect.DeepEqual(user, user1) {
		t.Errorf("expected %v users, got %v", user1, user)
	}
}

func TestPatchUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	insertTestUser(mockRepo, &user1)

	tests := []struct {
		name    string
		url     string
		body    map[string]any
		expCode int
	}{
		{
			name:    "success",
			url:     "/api/v1/users/1",
			body:    map[string]any{"full_name": "John Doe Jr."},
			expCode: http.StatusNoContent,
		},
		{
			name:    "user not found",
			url:     "/api/v1/users/2",
			expCode: http.StatusNotFound,
		},
		{
			name:    "bad request",
			url:     "/api/v1/users/1",
			body:    map[string]any{"email": ""},
			expCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := requester(http.MethodPatch, tt.url, tt.body, mockRepo)

			if rr.Code != tt.expCode {
				t.Errorf("expected status %d, got %d", tt.expCode, rr.Code)
			}

			if tt.expCode == http.StatusNoContent {
				rr = requester(http.MethodGet, "/api/v1/users/id/1", nil, mockRepo)

				var user model.User
				if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if user.FullName != "John Doe Jr." {
					t.Errorf("expected %v user full_name, got %v", "John Doe Jr.", user.FullName)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	insertTestUser(mockRepo, &user1)

	rr := requester(http.MethodDelete, "/api/v1/users/1", nil, mockRepo)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rr.Code)
	}

	if userExists(mockRepo, 1) {
		t.Errorf("user was not deleted from the database")
	}
}
