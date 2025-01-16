package entity

var Roles = struct {
	Guest      uint
	User       uint
	Authorized uint
}{
	Guest:      0,
	User:       1,
	Authorized: 3,
}
