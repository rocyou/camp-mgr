package dao

import (
	"camp-mgr/app/campmgr/internal/dao/mysql"
	"github.com/teamgram/marmota/pkg/stores/sqlx"
)

type Mysql struct {
	*sqlx.DB
	*mysql_dao.CampaignDAO
	*mysql_dao.RecipientDAO
	*mysql_dao.MessageDAO
}

func newMysqlDao(db *sqlx.DB) *Mysql {
	return &Mysql{
		DB:           db,
		CampaignDAO:  mysql_dao.NewCampaignDAO(db),
		RecipientDAO: mysql_dao.NewRecipientDAO(db),
		MessageDAO:   mysql_dao.NewMessageDAO(db),
	}
}
