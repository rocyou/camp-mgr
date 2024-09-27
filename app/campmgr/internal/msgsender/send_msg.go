package msgsender

import (
	"camp-mgr/app/campmgr/internal/config"
	"camp-mgr/app/campmgr/internal/dao/dataobject"
	mysql_dao "camp-mgr/app/campmgr/internal/dao/mysql"
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

func NewConsumer(svcCtx *svc.ServiceContext, conf kafka.KafkaConsumerConf) *kafka.BatchConsumerGroup {
	// Batch consumption: each batch consumes up to 1000 messages.
	// If fewer than 1000 messages are accumulated within 100ms, the available messages will still be consumed in one batch.
	s := kafka.MustKafkaBatchConsumer(&conf)
	s.RegisterHandlers(
		func(triggerID string, idList []string) {
		},

		func(value kafka.MsgChannelValue) {
			logx.Info("BatchConsumerGroup - AggregationID: ", value.AggregationID, ", TriggerID: ", value.TriggerID, ", Len: ", len(value.MsgList))
			//var wg sync.WaitGroup
			perMsgTableData := make(map[string][]*dataobject.Message)
			for _, msg := range value.MsgList {
				key := value.AggregationID
				value1 := msg.MsgData
				logx.Infof("key: %s, value: %s", key, value1)

				var message dataobject.Message
				_ = jsoniter.Unmarshal(msg.MsgData, &message)

				campaignId, _ := strconv.ParseInt(key, 10, 64)
				tableSharding := mysql_dao.CalMsgTable(campaignId, int64(svcCtx.Config.MsgTableShardingSize))

				logx.Infof("campaignId: %d, value: %+v", campaignId, message)
				if _, exists := perMsgTableData[tableSharding]; exists {
					perMsgTableData[tableSharding] = append(perMsgTableData[tableSharding], &message)
				} else {
					perMsgTableData[tableSharding] = []*dataobject.Message{&message}
				}
			}

			callWhatsappAPI := func(sendInfo *dataobject.Message) bool {
				success := true
				// WhatsApp API call logic here...
				time.Sleep(500 * time.Millisecond)
				if sendInfo.Id%2 == 0 {
					success = true
				} else {
					success = false
				}
				return success
			}

			for tableName, values := range perMsgTableData {

				// It may instantaneous increase the pressure on the database use go routine
				go func(table string, messages []*dataobject.Message) {
					var wg sync.WaitGroup
					sendResult := make(chan *dataobject.Message, len(messages))
					for _, msg := range messages {
						wg.Add(1)
						go func(msg1 *dataobject.Message) {
							defer wg.Done()
							//msg1.SendAt = time.Now().Unix()
							if callWhatsappAPI(msg1) {
								msg1.Status = 1 //send success
							} else {
								msg1.Status = 2 //send failed
							}

							logx.Infof("send result: %+v", msg1)
							sendResult <- msg1

						}(msg)
					}
					wg.Wait()
					close(sendResult)
					doList := make([]*dataobject.Message, 0, len(messages))
					for val := range sendResult {
						doList = append(doList, val)
					}
					// need add service monitor here when BatchUpdateById failed
					_, err := svcCtx.Dao.MessageDAO.BatchUpdateById(context.Background(), tableName, doList, time.Now().Unix())
					if err != nil {
						// need add service monitor here
						logx.Infof("update send result in db: %+v", err)
					}

				}(tableName, values)
			}
		})

	return s
}
