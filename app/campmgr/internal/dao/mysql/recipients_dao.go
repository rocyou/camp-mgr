package mysql_dao

import (
	"camp-mgr/app/campmgr/internal/dao/dataobject"
	"database/sql"
	"fmt"
	"github.com/teamgram/marmota/pkg/stores/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

var _ *sql.Result
var _ = fmt.Sprintf
var _ = strings.Join

type RecipientDAO struct {
	db *sqlx.DB
}

func NewRecipientDAO(db *sqlx.DB) *RecipientDAO {
	return &RecipientDAO{
		db: db,
	}
}

func (dao *RecipientDAO) InsertOrUpdateBulkTx(tx *sqlx.Tx, doList []*dataobject.Recipients) (lastInsertId, rowsAffected int64, err error) {
	if len(doList) == 0 {
		return
	}

	var (
		query = "INSERT INTO recipients (phone, name) VALUES "
		vals  []string
		args  []interface{}
	)

	for _, do := range doList {
		vals = append(vals, "(?, ?)")
		args = append(args, do.Phone, do.Name)
	}

	query += strings.Join(vals, ", ") + " ON DUPLICATE KEY UPDATE phone = VALUES(phone), name = VALUES(name)"

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

func (dao *RecipientDAO) DeleteTx(tx *sqlx.Tx, phone string) (rowsAffected int64, err error) {
	var (
		query   = "update recipients set deleted = 1 where phone = ? and deleted = 0"
		rResult sql.Result
	)
	rResult, err = tx.Exec(query, phone)

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
