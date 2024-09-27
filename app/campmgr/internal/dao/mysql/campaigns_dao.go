package mysql_dao

import (
	"camp-mgr/app/campmgr/internal/dao/dataobject"
	"context"
	"database/sql"
	"fmt"
	"github.com/teamgram/marmota/pkg/stores/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

var _ *sql.Result
var _ = fmt.Sprintf
var _ = strings.Join

type CampaignDAO struct {
	db *sqlx.DB
}

func NewCampaignDAO(db *sqlx.DB) *CampaignDAO {
	return &CampaignDAO{
		db: db,
	}
}

func (dao *CampaignDAO) InsertOrUpdateTx(tx *sqlx.Tx, do *dataobject.Campaign) (lastInsertId, rowsAffected int64, err error) {
	var (
		query = "insert into campaigns(campaign_name, campaign_id, message_template, scheduled_time, csv_file_path ) values (:campaign_name, :campaign_id, :message_template, :scheduled_time, :csv_file_path) on duplicate key update campaign_name = values(campaign_name), campaign_id = values(campaign_id), message_template = values(message_template), scheduled_time = values(scheduled_time), csv_file_path = values(csv_file_path)"
		r     sql.Result
	)

	r, err = tx.NamedExec(query, do)
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("namedExec in InsertOrUpdateTx(%v), error: %v", do, err)
		return
	}

	lastInsertId, err = r.LastInsertId()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("lastInsertId in InsertOrUpdateTx(%v)_error: %v", do, err)
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("rowsAffected in InsertOrUpdateTx(%v)_error: %v", do, err)
	}

	return
}

func (dao *CampaignDAO) SelectListByCampaignId(ctx context.Context, campaignId string) (rList []dataobject.Campaign, err error) {

	var (
		query  = "select id, campaign_name, campaign_id, message_template, scheduled_time, csv_file_path from campaigns where campaignId = ?"
		values []dataobject.Campaign
	)
	err = dao.db.QueryRowsPartial(ctx, &values, query, campaignId)

	if err != nil {
		logx.WithContext(ctx).Errorf("queryx in SelectListByCampaignId(_), error: %v", err)
		return
	}

	rList = values

	return
}

func (dao *CampaignDAO) InsertBulkTx(tx *sqlx.Tx, do []*dataobject.Campaign) (lastInsertId, rowsAffected int64, err error) {
	var (
		query = "insert into campaigns(campaign_name, campaign_id, message_template, scheduled_time, csv_file_path ) values (:campaign_name, :campaign_id, :message_template, :scheduled_time, :csv_file_path)"
		r     sql.Result
	)

	r, err = tx.NamedExec(query, do)
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("namedExec in Insert(%v), error: %v", do, err)
		return
	}

	lastInsertId, err = r.LastInsertId()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("lastInsertId in Insert(%v)_error: %v", do, err)
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("rowsAffected in Insert(%v)_error: %v", do, err)
	}

	return
}

func (dao *CampaignDAO) DeleteTx(tx *sqlx.Tx, campaignId int64) (rowsAffected int64, err error) {
	var (
		query   = "update campaigns set deleted = 1 where campaign_id = ? and deleted = 0"
		rResult sql.Result
	)
	rResult, err = tx.Exec(query, campaignId)

	if err != nil {
		logx.WithContext(tx.Context()).Errorf("exec in Delete(_), error: %v", err)
		return
	}

	rowsAffected, err = rResult.RowsAffected()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("rowsAffected in Delete(_), error: %v", err)
	}

	return
}
