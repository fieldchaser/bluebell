package models

import "time"

type Community struct {
	CommunityId   int64  `db:"community_id" json:"id"`
	CommunityName string `db:"community_name" json:"name"`
}

type CommunityDetail struct {
	CommunityId   int64     `db:"community_id" json:"id"`
	CommunityName string    `db:"community_name" json:"name"`
	Introduction  string    `db:"introduction" json:"introduction"`
	CreateTime    time.Time `db:"create_time" json:"create_time"`
}
