package job

type CJobClient interface {
	AddCampaignJob(campaignId, scheduleTime int64) error
	UpdateCampaignJob(campaignId, newScheduleTime int64) error
	RemoveCampaignJob(campaignId int64)
	Start()
	Stop()
}
