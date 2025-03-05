package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"web_framework/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlstr := `SELECT community_id, community_name FROM community`
	if err := db.Select(&communityList, sqlstr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("no community found")
			err = nil
		}
	}
	return
}

func GetCommunityDetailList(id int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlstr := `SELECT community_id, community_name, introduction, create_time FROM community WHERE community_id = ?`
	if err := db.Get(communityDetail, sqlstr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return communityDetail, nil
}
