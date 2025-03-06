package logic

import (
	"go.uber.org/zap"
	"web_framework/dao/mysql"
	"web_framework/dao/redis"
	"web_framework/models"
	"web_framework/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	//1.生成post_id
	p.PostID = snowflake.GenID()
	//2.存到数据库里
	err = mysql.CreatePost(p)
	if err != nil {
		zap.L().Error("mysql.CreatePost failed", zap.Error(err))
		return
	}
	//3.在redis里同步时间
	err = redis.CreatePost(p.PostID)
	return
}

func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	//1.根据uid获取post
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Error(err))
		return
	}
	//2.根据author_id获取author_name
	author_id := post.AuthorId
	Author, err := mysql.GetAuthorNameById(author_id)
	if err != nil {
		zap.L().Error("mysql.GetAuthorNameById(author_id) failed", zap.Error(err))
		return
	}
	//3.根据community_id获取community详情
	community, err := mysql.GetCommunityDetailList(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailList(post.CommunityID) failed", zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      Author.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

func GetPostDetail(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostDetail(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	var community *models.CommunityDetail

	for _, post := range posts {
		user_id := post.AuthorId
		user, err := mysql.GetAuthorNameById(user_id)
		if err != nil {
			zap.L().Error("mysql.GetAuthorNameById(author_id) failed", zap.Error(err))
			continue
		}
		//3.根据community_id获取community详情
		community, err = mysql.GetCommunityDetailList(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailList(post.CommunityID) failed", zap.Error(err))
			continue
		}
		info := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, info)
	}

	return
}
