package msgsyncer

import (
	"context"
	kafka "github.com/teamgram/marmota/pkg/mq"
	"github.com/zeromicro/go-zero/core/jsonx"
	"strconv"
)

type SyncMsg struct {
	Id          int64  `json:"id"`
	CampaignId  int64  `json:"campaign_id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	MessageData string `json:"message_data"`
	LastMessage bool   `json:"last_message"`
}

type DefaultSyncMsgClient struct {
	cli *kafka.Producer
}

func NewSyncMsgClient(cli *kafka.Producer) SyncMsgClient {
	return &DefaultSyncMsgClient{
		cli: cli,
	}
}

func (m *DefaultSyncMsgClient) SyncMessage(ctx context.Context, campaignId int64, in *SyncMsg) error {
	return m.syncMessage(
		ctx,
		campaignId,
		in)
}

func (m *DefaultSyncMsgClient) syncMessage(ctx context.Context, k int64, in interface{}) error {
	var (
		b   []byte
		err error
	)

	b, err = jsonx.Marshal(in)
	if err != nil {
		return err
	}
	_, _, err = m.cli.SendMessage(ctx, strconv.FormatInt(k, 10), b)
	if err != nil {
		return err
	}

	return nil
}
