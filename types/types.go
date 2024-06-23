package types

type User struct {
	UserID   string `json:"userid"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}
