package logic

import (
	"web_framework/dao/mysql"
	"web_framework/models"
	"web_framework/pkg/jwt"
	"web_framework/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//需要与dao层交互
	//1.判断用户是否存在
	var ok bool
	ok, err = mysql.QueryUserByName(p.Username)
	if err != nil {
		return err
	}
	if ok {
		return mysql.ErrorUserExist
	}
	//2.生成uid
	userID := snowflake.GenID()
	//3.创建用户对象并存进mysql
	user := &models.User{
		Username: p.Username,
		//这里的密码是要加密过后的
		Password: p.Password,
		UserId:   userID,
	}
	err = mysql.InsertUser(user)
	if err != nil {
		return err
	}
	return nil
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成jwt的token
	token, err := jwt.GenToken(user.UserId, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return user, nil
}
