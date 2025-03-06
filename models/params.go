package models

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamPostVote struct {
	//userid 可以从请求中直接获取
	PostId    string `json:"post_id" binding:"required"`
	Direction int8   `json:"direction,string" binding:"oneof=0 1 -1"`
}
