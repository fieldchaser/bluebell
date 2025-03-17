package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePostHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandlers)

	body := `{
		"community_id": 1,
		"title": "test",
		"content": "just a test"
	}`

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	//判断响应的内容是不是按预期返回了需要登录的错误

	//方法1：判断响应内容中是不是包含指定的字符串
	//assert.Contains(t, w.Body.String(), "需要登陆")

	//方法2：将响应的内容反序列化到ResponseData 然后判断字段与预期是否一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.body failed, err: %v", err)
	}
	assert.Equal(t, CodeNeedLogin, res.Code)
}
