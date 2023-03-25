package auth

type MakeAuthRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}
