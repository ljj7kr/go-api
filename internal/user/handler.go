package user

import (
	"context"
	"errors"
	"net/http"

	"go-api/internal/gen/openapi"
	"go-api/internal/httpx"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Handler struct {
	service ServiceInterface
}

type ServiceInterface interface {
	CreateUser(ctx context.Context, input CreateUserRequest) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	ListUsers(ctx context.Context, input ListUsersRequest) (ListUsersResult, error)
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request, params openapi.ListUsersParams) {
	result, err := h.service.ListUsers(r.Context(), ListUsersRequest{
		Page: valueOrDefault(params.Page, 1),
		Size: valueOrDefault(params.Size, 20),
	})
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, openapi.ListUsersResponse{
		HasNext: result.HasNext,
		Items:   toUserResponses(result.Items),
		Page:    result.Page,
		Size:    result.Size,
	})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req openapi.CreateUserRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "잘못된 요청입니다", nil)
		return
	}

	createdUser, err := h.service.CreateUser(r.Context(), CreateUserRequest{
		Name:  req.Name,
		Email: string(req.Email),
	})
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, toUserResponse(createdUser))
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request, id int64) {
	foundUser, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toUserResponse(foundUser))
}

func (h *Handler) writeServiceError(w http.ResponseWriter, err error) {
	var invalidInputErr *InvalidInputError
	switch {
	case errors.As(err, &invalidInputErr):
		httpx.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "잘못된 요청입니다", toErrorDetails(invalidInputErr.Details))
	case errors.Is(err, ErrDuplicateUserEmail):
		httpx.WriteError(w, http.StatusConflict, "USER_ALREADY_EXISTS", "이미 존재하는 사용자입니다", nil)
	case errors.Is(err, ErrUserNotFound):
		httpx.WriteError(w, http.StatusNotFound, "USER_NOT_FOUND", "사용자를 찾을 수 없습니다", nil)
	default:
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "서버 내부 오류입니다", nil)
	}
}

func toUserResponse(user User) openapi.UserResponse {
	return openapi.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		CreatedAt: user.CreatedAt,
	}
}

func toUserResponses(users []User) []openapi.UserResponse {
	items := make([]openapi.UserResponse, 0, len(users))
	for _, foundUser := range users {
		items = append(items, toUserResponse(foundUser))
	}

	return items
}

func toErrorDetails(details []FieldError) []openapi.ErrorDetail {
	if len(details) == 0 {
		return nil
	}

	result := make([]openapi.ErrorDetail, 0, len(details))
	for _, detail := range details {
		result = append(result, openapi.ErrorDetail{
			Field:  detail.Field,
			Reason: detail.Reason,
		})
	}

	return result
}

func valueOrDefault(value *int, defaultValue int) int {
	if value == nil {
		return defaultValue
	}

	return *value
}
