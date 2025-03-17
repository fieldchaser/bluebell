package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"testing"
	"web_framework/models"
)

//var db *sqlx.DB

func initt() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"127.0.0.1",
		3306,
		"bluebell_test",
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	return
}

func init() {
	if err := initt(); err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		Title:       "Post Title",
		Content:     "Post Content",
		CommunityID: 1,
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost failed, err: %v\n", err)
	}
	t.Logf("CreatePost: success")
}
