package msgsyncer

import (
	"context"
)

type SyncMsgClient interface {
	SyncMessage(ctx context.Context, campaignId int64, in *SyncMsg) error
}
