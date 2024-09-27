package dao

import (
	"camp-mgr/app/campmgr/internal/config"
	"github.com/teamgram/marmota/pkg/stores/sqlx"
)

// Dao dao.
type Dao struct {
	*Mysql
}

// New new a dao and return.
func New(c config.Config) (dao *Dao) {
	db := sqlx.NewMySQL(&c.Mysql)
	return &Dao{
		Mysql: newMysqlDao(db),
	}
}
