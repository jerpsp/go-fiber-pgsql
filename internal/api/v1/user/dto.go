package user

type UserCreateRequest struct {
	Email     string `json:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" form:"password" validate:"required"`
	FirstName string `json:"first_name" form:"first_name" validate:"required"`
	LastName  string `json:"last_name" form:"last_name" validate:"omitempty"`
	Role      string `json:"role" form:"role" validate:"omitempty,oneof=admin user moderator"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" form:"first_name" validate:"required"`
	LastName  string `json:"last_name" form:"last_name" validate:"omitempty"`
}

type UserRoleUpdateRequest struct {
	Role string `json:"role" form:"role" validate:"required,oneof=admin user moderator"`
}

type PaginationRequest struct {
	Page  int `json:"page" query:"page"`
	Limit int `json:"limit" query:"limit"`
}

type PaginatedResponse struct {
	Users      []User `json:"users"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalPages int    `json:"total_pages"`
}
