package redis

import (
	"github.com/go-redis/redis"
	"math"
	"time"
)

// PostVote 处理投票的业务逻辑

//赞成票 + 432， 反对票 - 432
// 基于投票的排序算法，可以参照 https://www.ruanyifeng.com/blog/archives.html

/*
	direction=1时，有两种情况：1。用户之前没有投过票，现在投赞成票，2.用户之前投的是反对票，现在改成赞成票
	direction=0时，有两种情况：1。用户之前投过赞成票，现在取消投票，2.用户之前投的是反对票，现在取消投票
	direction=-1时，有两种情况：1。用户之前没有投过票，现在投反对票，2.用户之前投的是赞成票，现在改成反对票

	对投票的限制：发帖时间开始算起一个星期内的帖子允许投票，超过一个星期的禁止投票功能。
	1.到期之后，把redis中存的赞成票数及反对票数存到mysql表中
	2.到期之后删除 KeyPostVotedZSetPrefix 这个key
*/

const (
	VoteExpiration = 7 * 24 * 3600
	ScorePerVote   = 432
)

func VoteForPost(UserId string, PostId string, Direction float64) error {
	//1.先检验投票的限制
	current_time := time.Now().Unix()
	otime := rdb.ZScore(GetRedisKey(KeyPostTimeZSet), PostId).Val()
	if float64(current_time)-otime > VoteExpiration {
		err := ErrorVoteExpired
		return err
	}
	//2.根据direction更新帖子分数
	key := GetRedisKey(KeyPostVotedZSetPrefix) + PostId
	ov := rdb.ZScore(key, UserId).Val()
	diff := math.Abs(ov - Direction)
	var op float64
	if ov == Direction {
		return ErrorRepeated
	}
	if Direction > ov {
		op = 1
	} else {
		op = -1
	}
	//这里也需要事务
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(GetRedisKey(KeyPostScoreZSet), op*diff*ScorePerVote, PostId)

	//3.更新用户投票记录
	if Direction == 0 {
		pipeline.ZRem(GetRedisKey(KeyPostVotedZSetPrefix+PostId), UserId)
	} else {
		pipeline.ZAdd(GetRedisKey(KeyPostVotedZSetPrefix+PostId), redis.Z{
			Score:  Direction,
			Member: UserId,
		})
	}

	_, err := pipeline.Exec()
	return err
}
