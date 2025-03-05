package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_framework/logic"
	"web_framework/models"
)

// CreatePostHandlers 处理创建帖子的handler
func CreatePostHandlers(c *gin.Context) {
	//1.参数校验
	post := new(models.Post)
	if err := c.ShouldBindJSON(post); err != nil {
		zap.L().Debug("bind post failed", zap.Error(err))
		zap.L().Error("bind post failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//1.5 拿到当前的登录用户id
	ID, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("login needed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	post.AuthorId = ID
	//2.生成帖子
	if err := logic.CreatePost(post); err != nil {
		zap.L().Error("logic.CreatePost(post) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandlers 处理查询帖子详情的handler
func GetPostDetailHandlers(c *gin.Context) {
	//1.获取url中的id
	id := c.Param("id")
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zap.L().Error("parse pid failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id从数据库中查数据
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandlers 处理帖子列表的handler
func GetPostListHandlers(c *gin.Context) {
	//1.获取分页参数
	page, size := GetPageInfo(c)
	//2.查数据库并返回
	data, err := logic.GetPostDetail(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostDetail(page,size) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}
