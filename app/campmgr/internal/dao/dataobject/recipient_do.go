package dataobject

type Recipients struct {
	Id      int64  `db:"id" json:"id"`
	Phone   string `db:"phone" json:"phone"`
	Name    string `db:"name" json:"name"`
	Deleted int    `db:"deleted" json:"deleted"`
}
