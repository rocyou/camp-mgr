package job

import (
	"camp-mgr/app/campmgr/internal/config"
	"camp-mgr/app/campmgr/internal/dao"
	"camp-mgr/app/campmgr/internal/msgsyncer"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type DefaultJobClient struct {
	ctx     context.Context
	conf    config.Config
	dao     *dao.Dao
	mu      sync.Mutex
	syncCli msgsyncer.SyncMsgClient
	cron    *cron.Cron
	jobIDs  map[int64]cron.EntryID
	logx.Logger
}

func NewJobClient(ctx context.Context, conf config.Config, dao *dao.Dao, syncCli msgsyncer.SyncMsgClient) CJobClient {
	return &DefaultJobClient{
		ctx:     ctx,
		conf:    conf,
		dao:     dao,
		syncCli: syncCli,
		cron:    cron.New(cron.WithSeconds()),
		jobIDs:  make(map[int64]cron.EntryID),
	}
}

func (s *DefaultJobClient) AddCampaignJob(campaignId, scheduleTime int64) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := time.Unix(scheduleTime, 0)
	cronSpec := fmt.Sprintf("%d %d %d %d %d ?", t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()))

	logx.Infof("Adding scheduled job for campaign %d, cron expression: %s", campaignId, cronSpec)

	if jobID, exists := s.jobIDs[campaignId]; exists {
		s.cron.Remove(jobID)
	}

	jobID, err := s.cron.AddFunc(cronSpec, func() {
		logx.Infof("Scheduled job for campaign %d executed", campaignId)

		messages, err := s.dao.MessageDAO.SelectListByCampaignId(s.ctx, campaignId, int64(s.conf.MsgTableShardingSize))
		if err != nil {
			logx.Errorf("Failed to SelectListByCampaignId for campaign %d: %v", campaignId, err)
			return
		}

		// Do we need to rate-limit message production to Kafka here?
		for index, msg := range messages {
			msg1 := &msgsyncer.SyncMsg{
				Id:          msg.Id,
				CampaignId:  msg.CampaignId,
				Name:        msg.Name,
				Phone:       msg.Phone,
				MessageData: msg.MessageData,
			}
			//mark as last message for consumer patch process use
			if index == len(messages)-1 {
				msg1.LastMessage = true
			}
			s.syncCli.SyncMessage(s.ctx, msg1.CampaignId, msg1)
		}
		s.RemoveCampaignJob(campaignId) // Remove the job after execution
	})

	if err != nil {
		logx.Errorf("Failed to add scheduled job for campaign %d: %v", campaignId, err)
		return
	}
	s.jobIDs[campaignId] = jobID

	return
}

func (s *DefaultJobClient) UpdateCampaignJob(campaignId, newScheduleTime int64) (err error) {
	logx.Infof("Updating scheduled job for campaign %d", campaignId)
	return s.AddCampaignJob(campaignId, newScheduleTime)
}

func (s *DefaultJobClient) RemoveCampaignJob(campaignId int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if jobID, exists := s.jobIDs[campaignId]; exists {
		s.cron.Remove(jobID)
		delete(s.jobIDs, campaignId)
		logx.Infof("Scheduled job for campaign %d has been removed", campaignId)
	} else {
		logx.Infof("No scheduled job found for campaign %d", campaignId)
	}
}

func (s *DefaultJobClient) Start() {
	logx.Infof("Scheduled job start")
	s.cron.Start()
}

func (s *DefaultJobClient) Stop() {
	s.cron.Stop()
}
