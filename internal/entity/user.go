package entity

type User struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	NickName    string `json:"nickname"`
	Password    string `json:"password"`
	ConfirmPass string `json:"cfmpsw"`
	HashPass    string
}
