package chat

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"null-links/chat_service/internal/logic/chat"
	"null-links/chat_service/internal/svc"
	"null-links/chat_service/internal/types"
)

func ChatHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatHistoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewChatHistoryLogic(r.Context(), svcCtx)
		resp, err := l.ChatHistory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
