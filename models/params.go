package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamPostVote 投票数据
type ParamPostVote struct {
	//userid 可以从请求中直接获取
	PostId    string `json:"post_id" binding:"required"`
	Direction int8   `json:"direction,string" binding:"oneof=0 1 -1"`
}

// ParamPostList 获取帖子列表query string 参数
type ParamPostList struct {
	CommunityId int64  `json:"community_id" form:"community_id"`
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
}

// ParamCommunityPostList 按社区获取帖子列表query string 参数
//type ParamCommunityPostList struct {
//	*ParamPostList
//	CommunityId int64 `json:"community_id" form:"community_id"`
//}
