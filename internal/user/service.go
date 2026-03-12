package user

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var ErrDuplicateUserEmail = errors.New("duplicate user email")

type Repository interface {
	CreateUser(ctx context.Context, input CreateUserInput) (int64, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	ListUsers(ctx context.Context, input ListUsersInput) ([]User, error)
}

type Service struct {
	repo     Repository
	validate *validator.Validate
}

func NewService(repo Repository) *Service {
	validate := validator.New()

	// validation error 에서 Go field name 대신 json tag name 을 노출
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			return field.Name
		}
		return name
	})

	return &Service{
		repo:     repo,
		validate: validate,
	}
}

func (s *Service) CreateUser(ctx context.Context, input CreateUserRequest) (User, error) {
	if err := s.validate.Struct(input); err != nil {
		return User{}, newInvalidInputError(err)
	}

	id, err := s.repo.CreateUser(ctx, CreateUserInput{
		Name:  input.Name,
		Email: input.Email,
	})
	if err != nil {
		return User{}, err
	}

	createdUser, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("get created user: %w", err)
	}

	return createdUser, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (User, error) {
	if id < 1 {
		return User{}, &InvalidInputError{
			Details: []FieldError{
				{
					Field:  "id",
					Reason: "must be greater than or equal to 1",
				},
			},
		}
	}

	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) ListUsers(ctx context.Context, input ListUsersRequest) (ListUsersResult, error) {
	if input.Page == 0 {
		input.Page = 1
	}
	if input.Size == 0 {
		input.Size = 20
	}

	if err := s.validate.Struct(input); err != nil {
		return ListUsersResult{}, newInvalidInputError(err)
	}

	// 다음 페이지 존재 여부 계산을 위해 size + 1 건 조회
	limit := input.Size + 1
	offset := (input.Page - 1) * input.Size
	if limit > math.MaxInt32 || offset > math.MaxInt32 {
		return ListUsersResult{}, &InvalidInputError{
			Details: []FieldError{
				{
					Field:  "page",
					Reason: "is too large",
				},
			},
		}
	}

	users, err := s.repo.ListUsers(ctx, ListUsersInput{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return ListUsersResult{}, err
	}

	hasNext := len(users) > int(input.Size)
	if hasNext {
		// 응답에는 요청한 크기만 포함
		users = users[:input.Size]
	}

	return ListUsersResult{
		Items:   users,
		Page:    input.Page,
		Size:    input.Size,
		HasNext: hasNext,
	}, nil
}

func newInvalidInputError(err error) error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	details := make([]FieldError, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		details = append(details, FieldError{
			Field:  validationError.Field(),
			Reason: validationMessage(validationError),
		})
	}

	return &InvalidInputError{Details: details}
}

func validationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must have at least %s character", err.Param())
	case "max":
		return fmt.Sprintf("must have at most %s characters", err.Param())
	default:
		return fmt.Sprintf("failed on %s validation", err.Tag())
	}
}
