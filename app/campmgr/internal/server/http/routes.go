package http

import (
	"camp-mgr/app/campmgr/internal/svc"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

func RegisterHandlers(engine *rest.Server, serverCtx *svc.ServiceContext) {
	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/v1/camp/create",
				Handler: createCamp(serverCtx),
			},
		},
	)
}
