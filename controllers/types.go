package controllers

type RegisterRequestData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
