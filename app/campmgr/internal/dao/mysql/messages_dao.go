package mysql_dao

import (
	"camp-mgr/app/campmgr/internal/dao/dataobject"
	"context"
	"database/sql"
	"fmt"
	"github.com/teamgram/marmota/pkg/stores/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"
	"time"
)

var _ *sql.Result
var _ = fmt.Sprintf
var _ = strings.Join

func CalMsgTable(campaignId, size int64) string {
	return fmt.Sprintf("messages_%s", strconv.FormatInt(campaignId%size, 10))
}

type MessageDAO struct {
	db *sqlx.DB
}

func NewMessageDAO(db *sqlx.DB) *MessageDAO {
	return &MessageDAO{
		db: db,
	}
}

func (dao *MessageDAO) InsertOrUpdateBulkTx(tx *sqlx.Tx, campaignId, sharding int64, doList []*dataobject.Message) (lastInsertId, rowsAffected int64, err error) {
	if len(doList) == 0 {
		return
	}

	var (
		//query = "INSERT INTO messages (campaign_id, name, phone, message_data, send_at, status) VALUES "
		query = fmt.Sprintf("INSERT INTO %s (campaign_id, name, phone, message_data, send_at, status) VALUES ", CalMsgTable(campaignId, sharding))

		vals []string
		args []interface{}
	)

	for _, do := range doList {
		vals = append(vals, "(?, ?, ?, ?, ?, ?)")
		args = append(args, do.CampaignId, do.Name, do.Phone, do.MessageData, do.SendAt, do.Status)
	}

	query += strings.Join(vals, ", ") + " ON DUPLICATE KEY UPDATE campaign_id = VALUES(campaign_id), name = VALUES(name), phone = VALUES(phone), message_data = VALUES(message_data), send_at = VALUES(send_at), status = VALUES(status)"

	r, err := tx.Exec(query, args...)
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("Exec in InsertOrUpdateBulk(%v), error: %v", doList, err)
		return
	}

	lastInsertId, err = r.LastInsertId()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("lastInsertId in InsertOrUpdateBulk(%v)_error: %v", doList, err)
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("rowsAffected in InsertOrUpdateBulk(%v)_error: %v", doList, err)
	}

	return
}

func (dao *MessageDAO) BatchUpdateById(ctx context.Context, tableName string, doList []*dataobject.Message, commonSendAt int64) (rowsAffected int64, err error) {
	if len(doList) == 0 {
		return
	}

	var idsStatus1 []string
	var idsStatus2 []string

	// Categorize ids into status = 1 and status = 2
	for _, do := range doList {
		if do.Status == 1 {
			idsStatus1 = append(idsStatus1, "?")
		} else if do.Status == 2 {
			idsStatus2 = append(idsStatus2, "?")
		}
	}

	args := []interface{}{commonSendAt}
	var rResult sql.Result

	/*
		UPDATE messages
		SET send_at = ?, status = 1
		WHERE id IN (?, ?, ?, ?);
	*/
	if len(idsStatus1) > 0 {
		queryStatus1 := fmt.Sprintf("UPDATE %s SET send_at = ?, status = 1 WHERE id IN (%s)", tableName, strings.Join(idsStatus1, ", "))
		for _, do := range doList {
			if do.Status == 1 {
				args = append(args, do.Id)
			}
		}

		rResult, err = dao.db.Exec(ctx, queryStatus1, args...)
		if err != nil {
			logx.WithContext(ctx).Errorf("Exec in BatchUpdateById (status = 1): error: %v", err)
			return
		}

		affected, _ := rResult.RowsAffected()
		rowsAffected += affected
	}

	/*
		UPDATE messages
		SET send_at = ?, status = 2
		WHERE id IN (?, ?, ?, ?);
	*/
	if len(idsStatus2) > 0 {
		args = []interface{}{commonSendAt}
		queryStatus2 := fmt.Sprintf("UPDATE %s SET send_at = ?, status = 2 WHERE id IN (%s)", tableName, strings.Join(idsStatus2, ", "))
		for _, do := range doList {
			if do.Status == 2 {
				args = append(args, do.Id)
			}
		}

		rResult, err = dao.db.Exec(ctx, queryStatus2, args...)
		if err != nil {
			logx.WithContext(ctx).Errorf("Exec in BatchUpdateById (status = 2): error: %v", err)
			return
		}

		affected, _ := rResult.RowsAffected()
		rowsAffected += affected
	}

	return
}

func (dao *MessageDAO) BatchUpdateByIdV2(ctx context.Context, db *sqlx.DB, tableName string, doList []*dataobject.Message) (err error) {
	if len(doList) == 0 {
		return
	}

	tR := sqlx.TxWrapper(ctx, db, func(tx *sqlx.Tx, result *sqlx.StoreResult) {
		// Create temporary table
		_, result.Err = db.Exec(ctx, "CREATE TEMPORARY TABLE temp_table (id BIGINT, send_at BIGINT, status INT);")
		if result.Err != nil {
			return
		}

		//// Batch insert data
		//queryInsert := "INSERT INTO temp_table (id, send_at, status) VALUES (?, ?, ?)"
		//for _, data := range doList {
		//	sendAtTime := time.Now().Unix()
		//	_, result.Err = db.Exec(ctx, queryInsert, data.Id, sendAtTime, data.Status)
		//	if result.Err != nil {
		//		return
		//	}
		//}

		// INSERT INTO temp_table (id, send_at, status) VALUES (1, 1727533161, 2), (2, 1727533161, 1) ...
		var values []string
		var args []interface{}

		for _, data := range doList {
			sendAtTime := time.Now().Unix()
			values = append(values, "(?, ?, ?)")
			args = append(args, data.Id, sendAtTime, data.Status)
		}

		queryInsert := fmt.Sprintf("INSERT INTO temp_table (id, send_at, status) VALUES %s", strings.Join(values, ", "))
		_, result.Err = db.Exec(ctx, queryInsert, args...)
		if result.Err != nil {
			return
		}

		/* UPDATE messages_xx m
		   JOIN temp_table t ON m.id = t.id
		   SET m.send_at = t.send_at, m.status = t.status;
		*/
		_, result.Err = db.Exec(ctx, fmt.Sprintf("UPDATE %s m JOIN temp_table t ON m.id = t.id SET m.send_at = t.send_at, m.status = t.status;", tableName))
		if result.Err != nil {
			return
		}

		// DROP TEMPORARY TABLE temp_table;
		_, result.Err = db.Exec(ctx, "DROP TEMPORARY TABLE temp_table;")
		if result.Err != nil {
			return
		}
	})

	return tR.Err
}

func (dao *MessageDAO) SelectListByCampaignId(ctx context.Context, campaignId, sharding int64) (rList []dataobject.Message, err error) {

	var (
		query = fmt.Sprintf("select id, campaign_id, name, phone, message_data, send_at, status from %s where campaign_id = %d and deleted = 0  ",
			CalMsgTable(campaignId, sharding), campaignId)
		values []dataobject.Message
	)
	err = dao.db.QueryRowsPartial(ctx, &values, query)

	if err != nil {
		logx.WithContext(ctx).Errorf("queryx in SelectListByCampaignId(_), error: %v", err)
		return
	}

	rList = values

	return
}

func (dao *MessageDAO) DeleteTx(tx *sqlx.Tx, campaignId int64) (rowsAffected int64, err error) {
	var (
		query   = "update messages set deleted = 1 where campaign_id = ? and deleted = 0"
		rResult sql.Result
	)
	rResult, err = tx.Exec(query, campaignId)

	if err != nil {
		logx.WithContext(tx.Context()).Errorf("exec in DeleteTx(_), error: %v", err)
		return
	}

	rowsAffected, err = rResult.RowsAffected()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("rowsAffected in DeleteTx(_), error: %v", err)
	}

	return
}
