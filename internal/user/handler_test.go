package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api/internal/gen/openapi"
)

type handlerServiceStub struct {
	listResult ListUsersResult
	listErr    error
}

func (s *handlerServiceStub) CreateUser(_ context.Context, _ CreateUserRequest) (User, error) {
	return User{}, errors.New("not implemented")
}

func (s *handlerServiceStub) GetUserByID(_ context.Context, _ int64) (User, error) {
	return User{}, errors.New("not implemented")
}

func (s *handlerServiceStub) UpdateUser(_ context.Context, _ int64, _ UpdateUserRequest) (User, error) {
	return User{}, errors.New("not implemented")
}

func (s *handlerServiceStub) ListUsers(_ context.Context, _ ListUsersRequest) (ListUsersResult, error) {
	if s.listErr != nil {
		return ListUsersResult{}, s.listErr
	}

	return s.listResult, nil
}

func (s *handlerServiceStub) DeleteUser(_ context.Context, _ int64) error {
	return errors.New("not implemented")
}

func Test사용자목록핸들러는_페이지응답을_반환한다(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&handlerServiceStub{
		listResult: ListUsersResult{
			Page:       1,
			Size:       2,
			TotalCount: 3,
			TotalPages: 2,
			HasNext:    true,
			HasPrev:    false,
			Items: []User{
				{ID: 11, Name: "열하나", Email: "11@example.com", CreatedAt: time.Unix(11, 0).UTC()},
				{ID: 10, Name: "열", Email: "10@example.com", CreatedAt: time.Unix(10, 0).UTC()},
			},
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/users?page=1&size=2", nil)
	response := httptest.NewRecorder()

	handler.ListUsers(response, request, openapi.ListUsersParams{
		Page: ptr(1),
		Size: ptr(2),
	})

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d", response.Code)
	}

	var result openapi.ListUsersResponse
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.Pagination.Page != 1 || result.Pagination.Size != 2 {
		t.Fatalf("unexpected paging info: page=%d size=%d", result.Pagination.Page, result.Pagination.Size)
	}
	if result.Pagination.TotalElements != 3 {
		t.Fatalf("unexpected total elements: got %d", result.Pagination.TotalElements)
	}
	if result.Pagination.TotalPages != 2 {
		t.Fatalf("unexpected total pages: got %d", result.Pagination.TotalPages)
	}
	if !result.Pagination.HasNext {
		t.Fatalf("unexpected hasNext: got false")
	}
	if result.Pagination.HasPrevious {
		t.Fatalf("unexpected hasPrevious: got true")
	}
	if len(result.Items) != 2 {
		t.Fatalf("unexpected item count: got %d", len(result.Items))
	}
	if result.Items[0].Email != "11@example.com" {
		t.Fatalf("unexpected first email: got %s", result.Items[0].Email)
	}
}

func Test사용자목록핸들러는_검증에러를_반환한다(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&handlerServiceStub{
		listErr: &InvalidInputError{
			Details: []FieldError{
				{Field: "size", Reason: "must have at most 100 characters"},
			},
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/users?page=1&size=101", nil)
	response := httptest.NewRecorder()

	handler.ListUsers(response, request, openapi.ListUsersParams{
		Page: ptr(1),
		Size: ptr(101),
	})

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d", response.Code)
	}

	var result openapi.ErrorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.Code != "INVALID_REQUEST" {
		t.Fatalf("unexpected error code: got %s", result.Code)
	}
	if result.Details == nil || len(*result.Details) != 1 {
		t.Fatal("expected error details")
	}
	if (*result.Details)[0].Field != "size" {
		t.Fatalf("unexpected field: got %s", (*result.Details)[0].Field)
	}
}

func ptr[T any](value T) *T {
	return &value
}
