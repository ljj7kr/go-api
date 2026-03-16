package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	gen "go-api/internal/gen/sqlc"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

var ErrUserNotFound = errors.New("user not found")

// User 는 repository 레이어에서 서비스로 전달하는 내부 모델이다
// sqlc generated type 를 그대로 바깥으로 노출하지 않기 위해 별도 정의한다
type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

// CreateUserInput 은 사용자 생성 입력값이다
type CreateUserInput struct {
	Name  string
	Email string
}

// UpdateUserInput 은 사용자 수정 입력값이다
type UpdateUserInput struct {
	ID    int64
	Name  string
	Email string
}

// ListUsersInput 은 목록 조회 입력값이다
type ListUsersInput struct {
	Limit  int32
	Offset int32
}

// UserRepository 는 사용자 영속성 접근을 담당한다
type UserRepository struct {
	q *gen.Queries
}

// NewUserRepository 는 sql.DB 기반 repository 를 생성한다
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		q: gen.New(db),
	}
}

// CreateUser 는 사용자를 생성하고 생성된 ID 를 반환한다
func (r *UserRepository) CreateUser(ctx context.Context, input CreateUserInput) (int64, error) {
	res, err := r.q.CreateUser(ctx, gen.CreateUserParams{
		Name:  input.Name,
		Email: input.Email,
	})
	if err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return 0, ErrDuplicateUserEmail
		}
		return 0, fmt.Errorf("create user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	return id, nil
}

// GetUserByID 는 ID 로 사용자를 조회한다
func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("get user by id: %w", err)
	}

	return toUser(row), nil
}

// GetUserByEmail 는 이메일로 사용자를 조회한다
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("get user by email: %w", err)
	}

	return toUser(row), nil
}

// UpdateUser 는 사용자를 수정한다
func (r *UserRepository) UpdateUser(ctx context.Context, input UpdateUserInput) error {
	res, err := r.q.UpdateUser(ctx, gen.UpdateUserParams{
		Name:  input.Name,
		Email: input.Email,
		ID:    input.ID,
	})
	if err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrDuplicateUserEmail
		}
		return fmt.Errorf("update user: %w", err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get update affected rows: %w", err)
	}
	if affectedRows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// ListUsers 는 사용자 목록을 조회한다
func (r *UserRepository) ListUsers(ctx context.Context, input ListUsersInput) ([]User, error) {
	rows, err := r.q.ListUsers(ctx, gen.ListUsersParams{
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	users := make([]User, 0, len(rows))
	for _, row := range rows {
		users = append(users, toUser(row))
	}

	return users, nil
}

// CountUsers 는 전체 사용자 수를 조회한다
func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	count, err := r.q.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}

// DeleteUser 는 사용자를 삭제한다
func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	res, err := r.q.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get delete affected rows: %w", err)
	}
	if affectedRows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func toUser(row gen.User) User {
	return User{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		CreatedAt: row.CreatedAt,
	}
}
