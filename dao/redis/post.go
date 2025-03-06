package redis

import (
	"github.com/go-redis/redis"
	"time"
)

func CreatePost(PostId int64) error {
	//注意这里帖子时间和帖子分数的初始化操作需要一起完成，所以需要事务
	pipeline := rdb.TxPipeline()
	//1.帖子时间
	pipeline.ZAdd(GetRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: PostId,
	})

	//2.帖子分数
	pipeline.ZAdd(GetRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: PostId,
	})

	_, err := pipeline.Exec()
	return err
}
