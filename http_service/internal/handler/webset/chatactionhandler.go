package webset

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nulltv/http_service/internal/logic/webset"
	"nulltv/http_service/internal/svc"
	"nulltv/http_service/internal/types"
)

func ChatActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := webset.NewChatActionLogic(r.Context(), svcCtx)
		resp, err := l.ChatAction(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
