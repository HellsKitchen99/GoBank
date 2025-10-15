package domain

type UserFromFront struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
