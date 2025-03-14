package redis

import (
	"errors"
)

var (
	ErrorVoteExpired = errors.New("已过投票时间")
	ErrorRepeated    = errors.New("不允许重复投相同票")
)
