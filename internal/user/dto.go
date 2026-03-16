package user

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Email string `json:"email" validate:"required,email,max=255"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Email string `json:"email" validate:"required,email,max=255"`
}

type ListUsersRequest struct {
	Page int `json:"page" validate:"min=1"`
	Size int `json:"size" validate:"min=1,max=100"`
}

type ListUsersResult struct {
	Items      []User
	Page       int
	Size       int
	TotalCount int64
	TotalPages int
	HasNext    bool
	HasPrev    bool
}

type FieldError struct {
	Field  string
	Reason string
}

type InvalidInputError struct {
	Details []FieldError
}

func (e *InvalidInputError) Error() string {
	return "invalid input"
}
