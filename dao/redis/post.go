package redis

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"web_framework/models"
)

func CreatePost(PostId, communityID int64) error {
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

	//把帖子id加到社区的set里面
	cKey := GetRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, PostId)
	_, err := pipeline.Exec()
	return err
}

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户传来的order参数来确定要查询的redis key
	key := GetRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = GetRedisKey(KeyPostScoreZSet)
	}
	//确定起始索引和终止索引
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostScore 根据id列表查询帖子的分数（赞成票的数量）
func GetPostScore(ids []string) (score []int64, err error) {
	//使用pipeline一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := GetRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return
	}
	score = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		score = append(score, int64(v))
	}
	return
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := GetRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyPostScoreZSet)
	}
	//使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	//针对新的zset就按照之前的逻辑去获取ids

	//社区的key
	cKey := GetRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityId)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityId))
	if rdb.Exists(key).Val() < 1 {
		//不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//从redis获取id
	//确定起始索引和终止索引
	return getIDsFromKey(key, p.Page, p.Size)
}
