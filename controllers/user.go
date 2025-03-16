package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"web_framework/dao/mysql"
	"web_framework/logic"
	"web_framework/models"
)

// SignUp 注册接口
// 实现注册功能的控制器
// @Summary 注册接口
// @Description 注册接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param object query models.ParamPostList false "查询参数"
// @Success 200 {object} _ResponsePostList
// @Router /signup [get]
func SignUp(c *gin.Context) {
	//1.参数的获取与校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("SignUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "参数校验错误",
	//	})
	//	return
	//}
	//2.业务逻辑
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

//实现登录功能的控制器

// Login 登录接口
// @Summary 登录接口
// @Description 登录接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param object query models.ParamLogin false "查询参数"
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func Login(c *gin.Context) {
	//1.验证请求参数
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务逻辑
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
		} else if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
		} else {
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	//3.返回响应
	ResponseSuccess(c, gin.H{
		"token":     fmt.Sprintf("%s", user.Token),
		"user_id":   user.UserId,
		"user_name": user.Username,
	})
}
