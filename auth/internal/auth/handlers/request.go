package handlers

type RegistredRequest struct {
	Name     string
	Password string
	Email    string
}

type LoginRequest struct {
	Email    string
	Password string
}
