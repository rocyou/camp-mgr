package http

import (
	"camp-mgr/app/campmgr/internal/svc"
	"github.com/zeromicro/go-zero/rest"
)

// New new a grpc server.
func New(ctx *svc.ServiceContext, c rest.RestConf) *rest.Server {
	srv := rest.MustNewServer(c)
	//srv.Use(BuildCommonParamsMiddleware())
	go func() {
		defer srv.Stop()
		RegisterHandlers(srv, ctx)
		srv.Start()
	}()
	return srv
}
