package dto

type LoginUserRequest struct {
	Login    string `json:"login" validate:"required,max=16"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"jwt_token"`
}

type RegisterUserRequest struct {
	Login    string `json:"login" validate:"required,max=16"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	Login string `json:"login"`
}

type UserResponse struct {
	Login string `json:"login"`
}
