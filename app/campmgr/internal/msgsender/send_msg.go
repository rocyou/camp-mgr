package msgsender

import (
	"camp-mgr/app/campmgr/internal/config"
	"camp-mgr/app/campmgr/internal/dao/dataobject"
	mysql_dao "camp-mgr/app/campmgr/internal/dao/mysql"
	"camp-mgr/app/campmgr/internal/msgsyncer"
	"camp-mgr/app/campmgr/internal/svc"
	"context"
	jsoniter "github.com/json-iterator/go"
	kafka "github.com/teamgram/marmota/pkg/mq"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"sync"
	"time"
)

type (
	Config = config.Config
)

func NewConsumerV2(svcCtx *svc.ServiceContext, conf kafka.KafkaConsumerConf) *kafka.BatchConsumerGroup {
	// Batch consumption: each batch consumes up to 1000 messages.
	// If fewer than 1000 messages are accumulated within 100ms, the available messages will still be consumed in one batch.
	var syncMsgMap = sync.Map{}
	s := kafka.MustKafkaBatchConsumer(&conf)
	s.RegisterHandlers(
		func(triggerID string, idList []string) {
		},
		func(value kafka.MsgChannelValue) {
			logx.Info("BatchConsumerGroup - AggregationID: ", value.AggregationID, ", TriggerID: ", value.TriggerID, ", Len: ", len(value.MsgList))
			for _, msg := range value.MsgList {
				key := value.AggregationID
				value1 := msg.MsgData
				logx.Infof("key: %s, value: %s", key, value1)

				var message msgsyncer.SyncMsg
				_ = jsoniter.Unmarshal(msg.MsgData, &message)

				if !message.LastMessage {
					if value, ok := syncMsgMap.Load(key); ok {
						msgList := value.([]*msgsyncer.SyncMsg)
						msgList = append(msgList, &message)
						syncMsgMap.Store(key, msgList)
					} else {
						syncMsgMap.Store(key, []*msgsyncer.SyncMsg{&message})
					}
				} else {
					// if its last message
					// first store the last message of this campaign
					msgList := []*msgsyncer.SyncMsg{}
					if value, ok := syncMsgMap.Load(key); ok {
						msgList = value.([]*msgsyncer.SyncMsg)
						msgList = append(msgList, &message)
						syncMsgMap.Store(key, msgList)
					} else {
						syncMsgMap.Store(key, []*msgsyncer.SyncMsg{&message})
					}
					// then patch process this campaign in one go routine
					go func(msgList []*msgsyncer.SyncMsg, batchSize int) {
						// Split msgList into batches based on the batch size
						for i := 0; i < len(msgList); i += batchSize {
							end := i + batchSize
							if end > len(msgList) {
								end = len(msgList)
							}

							var wg sync.WaitGroup
							sendResult := make(chan *dataobject.Message, len(msgList))
							// Process each batch of messages
							batch := msgList[i:end]
							for _, m := range batch {
								wg.Add(1)
								go func(msg *msgsyncer.SyncMsg) {
									defer wg.Done()
									msgResult := &dataobject.Message{
										Id:          msg.Id,
										CampaignId:  msg.CampaignId,
										Name:        msg.Name,
										Phone:       msg.Phone,
										MessageData: msg.MessageData,
									}

									if CallWhatsappAPI(msg) {
										msgResult.SendAt = time.Now().Unix()
										msgResult.Status = 1 // send success
									} else {
										msgResult.SendAt = time.Now().Unix()
										msgResult.Status = 2 // send failed
									}

									logx.Infof("send result: %+v", msgResult)
									sendResult <- msgResult
								}(m)
							}
							// Wait for all messages in the current batch to be processed
							wg.Wait()
							close(sendResult)

							doList := make([]*dataobject.Message, 0, len(sendResult))
							for val := range sendResult {
								doList = append(doList, val)
							}
							campaignId, _ := strconv.ParseInt(key, 10, 64)
							tableSharding := mysql_dao.CalMsgTable(campaignId, int64(svcCtx.Config.MsgTableShardingSize))
							// Persist all results into the database for this campaign
							err := svcCtx.Dao.MessageDAO.BatchUpdateByIdV2(context.Background(), svcCtx.Dao.DB, tableSharding, doList)
							if err != nil {
								// need add service monitor here
								logx.Errorf("update send result in db: %+v", err)
							}
						}
					}(msgList, 1000) // Process 1000 messages per batch
				}
			}
		})

	return s
}

func CallWhatsappAPI(sendInfo *msgsyncer.SyncMsg) bool {
	success := true
	// WhatsApp API call logic here...
	time.Sleep(10 * time.Millisecond)
	if sendInfo.Id%2 == 0 {
		success = true
	} else {
		success = false
	}
	return success
}
