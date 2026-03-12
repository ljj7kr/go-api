package user

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"
)

type listUsersRepoStub struct {
	users     []User
	listInput ListUsersInput
	listErr   error
}

func (s *listUsersRepoStub) CreateUser(_ context.Context, _ CreateUserInput) (int64, error) {
	return 0, errors.New("not implemented")
}

func (s *listUsersRepoStub) GetUserByID(_ context.Context, _ int64) (User, error) {
	return User{}, errors.New("not implemented")
}

func (s *listUsersRepoStub) ListUsers(_ context.Context, input ListUsersInput) ([]User, error) {
	s.listInput = input
	if s.listErr != nil {
		return nil, s.listErr
	}

	return s.users, nil
}

func Test사용자목록조회시_기본페이지네이션을_적용한다(t *testing.T) {
	t.Parallel()

	repo := &listUsersRepoStub{
		users: []User{
			{ID: 3, Name: "세번째", Email: "3@example.com", CreatedAt: time.Unix(3, 0)},
			{ID: 2, Name: "두번째", Email: "2@example.com", CreatedAt: time.Unix(2, 0)},
		},
	}
	service := NewService(repo)

	result, err := service.ListUsers(context.Background(), ListUsersRequest{})
	if err != nil {
		t.Fatalf("list users failed: %v", err)
	}

	if result.Page != 1 {
		t.Fatalf("unexpected page: got %d", result.Page)
	}
	if result.Size != 20 {
		t.Fatalf("unexpected size: got %d", result.Size)
	}
	if result.HasNext {
		t.Fatalf("unexpected hasNext: got true")
	}
	if repo.listInput.Limit != 21 {
		t.Fatalf("unexpected limit: got %d", repo.listInput.Limit)
	}
	if repo.listInput.Offset != 0 {
		t.Fatalf("unexpected offset: got %d", repo.listInput.Offset)
	}
	if len(result.Items) != 2 {
		t.Fatalf("unexpected item count: got %d", len(result.Items))
	}
}

func Test사용자목록조회시_다음페이지여부를_계산한다(t *testing.T) {
	t.Parallel()

	repo := &listUsersRepoStub{
		users: []User{
			{ID: 9, Name: "아홉", Email: "9@example.com", CreatedAt: time.Unix(9, 0)},
			{ID: 8, Name: "여덟", Email: "8@example.com", CreatedAt: time.Unix(8, 0)},
			{ID: 7, Name: "일곱", Email: "7@example.com", CreatedAt: time.Unix(7, 0)},
		},
	}
	service := NewService(repo)

	result, err := service.ListUsers(context.Background(), ListUsersRequest{
		Page: 2,
		Size: 2,
	})
	if err != nil {
		t.Fatalf("list users failed: %v", err)
	}

	if !result.HasNext {
		t.Fatalf("unexpected hasNext: got false")
	}
	if repo.listInput.Limit != 3 {
		t.Fatalf("unexpected limit: got %d", repo.listInput.Limit)
	}
	if repo.listInput.Offset != 2 {
		t.Fatalf("unexpected offset: got %d", repo.listInput.Offset)
	}
	if len(result.Items) != 2 {
		t.Fatalf("unexpected item count: got %d", len(result.Items))
	}
	if result.Items[1].ID != 8 {
		t.Fatalf("unexpected last user id: got %d", result.Items[1].ID)
	}
}

func Test사용자목록조회시_잘못된크기를_검증한다(t *testing.T) {
	t.Parallel()

	service := NewService(&listUsersRepoStub{})

	_, err := service.ListUsers(context.Background(), ListUsersRequest{
		Page: 1,
		Size: 101,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	var inputErr *InvalidInputError
	if !errors.As(err, &inputErr) {
		t.Fatalf("expected InvalidInputError: %T", err)
	}
	if len(inputErr.Details) != 1 {
		t.Fatalf("unexpected detail count: got %d", len(inputErr.Details))
	}
	if inputErr.Details[0].Field != "size" {
		t.Fatalf("unexpected field: got %s", inputErr.Details[0].Field)
	}
}

func Test사용자목록조회시_너무큰페이지를_검증한다(t *testing.T) {
	t.Parallel()

	service := NewService(&listUsersRepoStub{})

	_, err := service.ListUsers(context.Background(), ListUsersRequest{
		Page: math.MaxInt32,
		Size: 100,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	var inputErr *InvalidInputError
	if !errors.As(err, &inputErr) {
		t.Fatalf("expected InvalidInputError: %T", err)
	}
	if len(inputErr.Details) != 1 {
		t.Fatalf("unexpected detail count: got %d", len(inputErr.Details))
	}
	if inputErr.Details[0].Field != "page" {
		t.Fatalf("unexpected field: got %s", inputErr.Details[0].Field)
	}
}
