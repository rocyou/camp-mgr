package config

import (
	kafka "github.com/teamgram/marmota/pkg/mq"
	"github.com/teamgram/marmota/pkg/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql      sqlx.Config
	SyncClient *kafka.KafkaProducerConf
	SendClient kafka.KafkaConsumerConf

	MsgTableShardingSize int `json:",default=8"`
}
