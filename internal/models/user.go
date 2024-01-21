package models

type User struct {
	ID       string
	Email    string
	UserName string
	Password string
	Admin    bool
}

type UpdateUserDTO struct {
	Email    *string
	UserName *string
	Password *string
}
