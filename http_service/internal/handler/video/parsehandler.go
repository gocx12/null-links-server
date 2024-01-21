package video

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nulltv/http_service/internal/logic/video"
	"nulltv/http_service/internal/svc"
	"nulltv/http_service/internal/types"
)

func ParseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ParseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := video.NewParseLogic(r.Context(), svcCtx)
		resp, err := l.Parse(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
