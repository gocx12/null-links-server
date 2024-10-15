package chat

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"null-links/http_service/internal/logic/chat"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
)

func HistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatHistoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewHistoryLogic(r.Context(), svcCtx)
		resp, err := l.History(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
