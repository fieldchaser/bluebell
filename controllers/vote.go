package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strconv"
	"web_framework/logic"
	"web_framework/models"
)

// PostVoteHandlers 处理投票功能的handler
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func PostVoteHandlers(c *gin.Context) {
	//1.参数校验
	p := new(models.ParamPostVote)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		} else if ok {
			errdata := removeTopStruct(errs.Translate(trans))
			ResponseErrorWithMsg(c, CodeInvalidParam, errdata)
			return
		}
	}
	//2.业务逻辑
	UserId, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	if err := logic.PostVote(strconv.Itoa(int(UserId)), p.PostId, float64(p.Direction)); err != nil {
		zap.L().Error("logic.PostVote failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}
