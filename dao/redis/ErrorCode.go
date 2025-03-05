package redis

import (
	"errors"
)

var (
	ErrorVoteExpired = errors.New("已过投票时间")
)
