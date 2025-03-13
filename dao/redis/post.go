package redis

import (
	"github.com/go-redis/redis"
	"time"
	"web_framework/models"
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

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户传来的order参数来确定要查询的redis key
	key := GetRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = GetRedisKey(KeyPostScoreZSet)
	}
	//确定起始索引和终止索引
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return rdb.ZRevRange(key, start, end).Result()
}
