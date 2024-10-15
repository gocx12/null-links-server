package chat

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"null-links/http_service/internal/logic/chat"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
)

func GetAllTopicHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatGetAllTopicReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewGetAllTopicLogic(r.Context(), svcCtx)
		resp, err := l.GetAllTopic(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
