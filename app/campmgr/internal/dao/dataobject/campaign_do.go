package dataobject

type Campaign struct {
	Id              int64  `db:"id" json:"id"`
	CampaignName    string `db:"campaign_name" json:"campaign_name"`
	CampaignId      int64  `db:"campaign_id" json:"campaign_id"`
	MessageTemplate string `db:"message_template" json:"message_template"`
	ScheduledTime   int64  `db:"scheduled_time" json:"scheduled_time"`
	CSVFilePath     string `db:"csv_file_path" json:"csv_file_path"`
	Deleted         int    `db:"deleted" json:"deleted"`
}
