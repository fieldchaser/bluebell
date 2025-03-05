package models

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
	UserId   int64  `db:"user_id"`
	Token    string `db:"token"`
}
