package mysql

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
	"web_framework/models"
)

func CreatePost(p *models.Post) (err error) {
	//1.sql语句
	sqlstr := `insert into post (post_id, title, content, community_id, author_id) values (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlstr, p.PostID, p.Title, p.Content, p.CommunityID, p.AuthorId)
	//2.返回响应
	return
}

func GetPostById(pid int64) (p *models.Post, err error) {
	sqlstr := `select post_id, title, content, author_id, community_id, create_time from post where post_id = ?`
	p = new(models.Post)
	err = db.Get(p, sqlstr, pid)
	return
}

func GetAuthorNameById(author_id int64) (author *models.User, err error) {
	Author := new(models.User)
	sqlstr := `select username from user where user_id = ?`
	if err := db.Get(Author, sqlstr, author_id); err != nil {
		zap.L().Error("db.GetAuthorNameById(author_id) failed", zap.Error(err))
		return nil, err
	}
	return Author, nil
}

func GetPostDetail(page, size int64) (data []*models.Post, err error) {
	sqlstr := `select post_id, title, content, author_id, community_id, create_time from post 
	limit ?,?
	`
	data = make([]*models.Post, 0, 2)
	err = db.Select(&data, sqlstr, (page-1)*size, size)
	return
}

func GetPostListByIDs(ids []string) (PostList []*models.Post, err error) {
	sqlstr := `select post_id, title, content, author_id, community_id, create_time from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)
	`
	query, args, err := sqlx.In(sqlstr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	err = db.Select(&PostList, query, args...)
	return
}
