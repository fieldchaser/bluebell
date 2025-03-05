package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"web_framework/models"
)

const secret = "fieldchaser"

func QueryUserByName(username string) (exist bool, err error) {
	var count int
	sqlstr := `select count(user_id) from user where username=?`
	err = db.Get(&count, sqlstr, username)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func InsertUser(user *models.User) (err error) {
	sqlstr := `insert into user (username, password, user_id) values (?, ?, ?)`
	user.Password = encryptPassWord(user.Password)
	_, err = db.Exec(sqlstr, user.Username, user.Password, user.UserId)
	if err != nil {
		return err
	}
	return nil
}

func encryptPassWord(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password
	sqlstr := `select user_id, username, password from user where username=?`
	err = db.Get(user, sqlstr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}

	password := encryptPassWord(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}
