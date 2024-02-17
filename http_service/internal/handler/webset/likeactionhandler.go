package webset

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"null-links/http_service/internal/logic/webset"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
)

func LikeActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LikeActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := webset.NewLikeActionLogic(r.Context(), svcCtx)
		resp, err := l.LikeAction(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
