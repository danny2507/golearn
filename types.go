package main

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	UserID      int     `json:"user_id"`
}
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type LoginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterRequestData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
