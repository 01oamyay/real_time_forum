package entity

type User struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ConfirmPass string `json:"cfmpsw"`
	HashPass    string
}
