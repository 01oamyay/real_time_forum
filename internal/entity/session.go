package entity

// Session models the persisted token tied to a user.
type Session struct {
	UserID uint
	Token  string
}
