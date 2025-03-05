package logic

import (
	"go.uber.org/zap"
	"web_framework/dao/redis"
)

// PostVote 处理投票的业务逻辑

// 基于投票的排序算法，可以参照 https://www.ruanyifeng.com/blog/archives.html

/*
	direction=1时，有两种情况：1。用户之前没有投过票，现在投赞成票，2.用户之前投的是反对票，现在改成赞成票
	direction=0时，有两种情况：1。用户之前投过赞成票，现在取消投票，2.用户之前投的是反对票，现在取消投票
	direction=-1时，有两种情况：1。用户之前没有投过票，现在投反对票，2.用户之前投的是赞成票，现在改成反对票

	对投票的限制：发帖时间开始算起一个星期内的帖子允许投票，超过一个星期的禁止投票功能。
	1.到期之后，把redis中存的赞成票数及反对票数存到mysql表中
	2.到期之后删除 KeyPostVotedZSetPrefix 这个key
*/

// PostVote 为帖子投票的函数
func PostVote(UserId string, PostId string, Direction float64) error {
	zap.L().Debug("PostVote", zap.String("UserId", UserId),
		zap.String("PostId", PostId),
		zap.Float64("Direction", Direction))
	return redis.VoteForPost(UserId, PostId, Direction)
}
