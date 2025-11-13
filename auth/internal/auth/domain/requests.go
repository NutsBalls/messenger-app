package domain

type RegisterRequest struct {
	Name     string
	Password string
	Email    string
}

type LoginRequest struct {
	Email    string
	Password string
}
