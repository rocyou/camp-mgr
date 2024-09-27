package server

import (
	"camp-mgr/app/campmgr/internal/config"
	"camp-mgr/app/campmgr/internal/msgsender"
	"camp-mgr/app/campmgr/internal/server/http"
	"camp-mgr/app/campmgr/internal/svc"
	"flag"
	"github.com/common-nighthawk/go-figure"
	kafka "github.com/teamgram/marmota/pkg/mq"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var (
	configFile   = flag.String("c", "../etc/campmgr.yaml", "the config file")
	consumerMode = flag.Bool("w", false, "run as kafka consumer")
)

type Server struct {
	httpSrv *rest.Server
	mq      *kafka.BatchConsumerGroup
}

func New() *Server {
	return new(Server)
}

func (s *Server) Initialize() error {
	figure.NewFigure("Msg Producer Start", "", true).Print()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.Infov(c)
	svcCtx := svc.NewServiceContext(c)
	svcCtx.Start()
	s.httpSrv = http.New(svcCtx, c.RestConf)

	if *consumerMode {
		s.mq = msgsender.NewConsumer(svcCtx, c.SendClient)
		go func() {
			s.mq.Start()
		}()
		figure.NewFigure("Msg Consumer Start", "", true).Print()
	} else {
		figure.NewFigure("Msg Producer Start", "", true).Print()
	}

	return nil
}

func (s *Server) RunLoop() {
}

func (s *Server) Destroy() {
}
