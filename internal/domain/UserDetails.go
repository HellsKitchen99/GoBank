package domain

type UserDetails struct {
	Id    int64    `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}
