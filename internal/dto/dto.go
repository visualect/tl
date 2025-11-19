package dto

type LoginUserRequest struct {
	Login    string `json:"login" validate:"required,max=16"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserRequest struct {
	Login    string `json:"login" validate:"required,max=16"`
	Password string `json:"password" validate:"required"`
}
