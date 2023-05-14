package models

type User struct {
	ID       string `db:"id"`
	Password string `db:"password"`
}
