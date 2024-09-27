package core

import (
	"camp-mgr/app/campmgr/internal/dao/dataobject"
	"camp-mgr/app/campmgr/internal/svc"
	"context"
	"encoding/csv"
	"github.com/teamgram/marmota/pkg/stores/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type CampReq struct {
	CampaignName    string `json:"campaign_name"`
	CampaignId      int64  `json:"campaign_id"`
	MessageTemplate string `json:"message_template"`
	ScheduledTime   int64  `json:"scheduled_time"`
	CSVFilePath     string ` json:"csv_file_path"`
}

type Camp struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCamp(ctx context.Context, svcCtx *svc.ServiceContext) *Camp {
	return &Camp{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (c Camp) AddCamp(req *CampReq) (err error) {
	file, err := os.Open(req.CSVFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return
	}
	messages := make([]*dataobject.Message, 0, len(records))
	recipients := make([]*dataobject.Recipients, 0, len(records))
	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		recipients = append(recipients, &dataobject.Recipients{
			Name:  record[0],
			Phone: record[1],
		})

		messages = append(messages, &dataobject.Message{
			CampaignId:  req.CampaignId,
			Name:        record[0],
			Phone:       record[1],
			MessageData: req.MessageTemplate,
			Status:      0,
		})
	}

	tR := sqlx.TxWrapper(c.ctx, c.svcCtx.DB, func(tx *sqlx.Tx, result *sqlx.StoreResult) {
		_, _, result.Err = c.svcCtx.Dao.CampaignDAO.InsertOrUpdateTx(tx, &dataobject.Campaign{
			CampaignName:    req.CampaignName,
			CampaignId:      req.CampaignId,
			MessageTemplate: req.MessageTemplate,
			ScheduledTime:   req.ScheduledTime,
			CSVFilePath:     req.CSVFilePath,
		})
		if result.Err != nil {
			return
		}

		_, _, result.Err = c.svcCtx.Dao.RecipientDAO.InsertOrUpdateBulkTx(tx, recipients)
		if result.Err != nil {
			return
		}

		shardingSize := int64(c.svcCtx.Config.MsgTableShardingSize)
		_, _, result.Err = c.svcCtx.Dao.MessageDAO.InsertOrUpdateBulkTx(tx, req.CampaignId, shardingSize, messages)
		if result.Err != nil {
			return
		}

	})

	if tR.Err != nil {
		return tR.Err
	}

	logx.WithContext(c.ctx).Infof("Create campaign successfully, campaignId %d", req.CampaignId)

	if err = c.svcCtx.CJobClient.AddCampaignJob(req.CampaignId, req.ScheduledTime); err != nil {
		return err
	}

	return
}
