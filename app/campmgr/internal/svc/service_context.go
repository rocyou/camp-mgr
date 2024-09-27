package svc

import (
	"camp-mgr/app/campmgr/internal/config"
	"camp-mgr/app/campmgr/internal/dao"
	"camp-mgr/app/campmgr/internal/job"
	"camp-mgr/app/campmgr/internal/msgsyncer"
	"context"
	kafka "github.com/teamgram/marmota/pkg/mq"
)

type ServiceContext struct {
	Config config.Config
	*dao.Dao
	msgsyncer.SyncMsgClient
	job.CJobClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	dao := dao.New(c)
	syncMsgClient := msgsyncer.NewSyncMsgClient(kafka.MustKafkaProducer(c.SyncClient))
	return &ServiceContext{
		Config:        c,
		Dao:           dao,
		CJobClient:    job.NewJobClient(context.Background(), c, dao, syncMsgClient),
		SyncMsgClient: syncMsgClient,
	}
}
