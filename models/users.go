package models

// User is struct
type User struct {
	ID    string `validate:"required,printascii,max=32"`
	Name  string `validate:"required,max=32"`
	Email string `validate:"required,email"`
}
