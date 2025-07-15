package user

type UserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"omitempty"`
	Role      string `json:"role" validate:"omitempty,oneof=admin user moderator"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name" validate:"omitempty"`
}

type UserRoleUpdateRequest struct {
	Role string `json:"role" validate:"required,oneof=admin user moderator"`
}
