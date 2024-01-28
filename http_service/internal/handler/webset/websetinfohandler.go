package webset

import (
	"net/http"

	"nulltv/http_service/internal/logic/webset"
	"nulltv/http_service/internal/svc"
	"nulltv/http_service/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func WebsetInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WebsetInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := webset.NewWebsetInfoLogic(r.Context(), svcCtx)
		resp, err := l.WebsetInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
