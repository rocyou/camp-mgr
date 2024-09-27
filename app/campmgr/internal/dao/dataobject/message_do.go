package dataobject

type Message struct {
	Id          int64  `db:"id" json:"id"`
	CampaignId  int64  `db:"campaign_id" json:"campaign_id"`
	Name        string `db:"name" json:"name"`
	Phone       string `db:"phone" json:"phone"`
	MessageData string `db:"message_data" json:"message_data"`
	SendAt      int64  `db:"send_at" json:"send_at"`
	Status      int    `db:"status" json:"status"`
	Deleted     int    `db:"deleted" json:"deleted"`
}
