package entity

// User represents a registered forum user profile.
type User struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	NickName    string `json:"nickname"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Password    string `json:"password"`
	ConfirmPass string `json:"cfmpsw"`
}

// UserInput captures the credentials submitted during sign in.
type UserInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
