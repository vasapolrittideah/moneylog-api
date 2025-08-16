package payload

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpRequest struct {
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required"`
	FullName string `json:"full_name" validate:"required"`
}

type SignUpResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
