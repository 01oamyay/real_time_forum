package entity

// Roles enumerates the access levels enforced by the HTTP handlers.
var Roles = struct {
	Guest      uint
	User       uint
	Authorized uint
}{
	Guest:      0,
	User:       1,
	Authorized: 3,
}
