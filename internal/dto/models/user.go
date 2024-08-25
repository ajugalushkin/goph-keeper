package models

type User struct {
	ID           int64  `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash []byte `json:"password" db:"password"`
}
