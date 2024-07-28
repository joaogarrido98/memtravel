package types

// ServerResponse holds the generic type for all responses in
type ServerResponse struct {
	Response string      `json:"rs"`
	Message  string      `json:"ms"`
	Data     interface{} `json:"dt,omitempty"`
}

type User struct {
	UserID   string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Active   bool   `json:"active,omitempty"`
}

type ChangePassword struct {
	OldPassword string `json:"op,omitempty"`
	NewPassword string `json:"np,omitempty"`
}

type RecoverPasswordTemplate struct {
	Password string
}
