package auth

type CreateUserRequest struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CreateUserResponse struct {
	AccessToken string `json:"access_token"`
}
