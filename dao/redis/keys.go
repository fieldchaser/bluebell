package redis

//redis key

//定义redis key的时候要注意用命名空间的形式来命名（用冒号分割），方便后续的查询及拆分

const (
	KeyPostPrefix          = "bluebell:"
	KeyPostTimeZSet        = "post:time"   //zset; 帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  //zset; 帖子及投票的分数
	KeyPostVotedZSetPrefix = "post:voted:" //zset; 记录用户及投票类型，参数是post_id
)

func GetRedisKey(key string) string {
	return KeyPostPrefix + key
}
