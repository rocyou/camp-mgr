package http

import (
	"camp-mgr/app/campmgr/internal/core"
	"camp-mgr/app/campmgr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func createCamp(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Debugf("request: %#v", r)

		type res struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		var err error

		req := new(core.CampReq)
		if err = httpx.ParseJsonBody(r, req); err != nil {
			httpx.WriteJson(w, http.StatusOK, res{Code: -1, Msg: "failed"})
			return
		}

		c := core.NewCamp(r.Context(), ctx)
		if err = c.AddCamp(req); err != nil {
			httpx.WriteJson(w, http.StatusOK, res{Code: -1, Msg: "failed"})
			return
		}

		httpx.WriteJson(w, http.StatusOK, res{
			Code: 0,
			Msg:  "success",
		})

	}
}
