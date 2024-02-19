package chat

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"null-links/chat_service/internal/logic/chat"
	"null-links/chat_service/internal/svc"
	"null-links/chat_service/internal/types"
)

func ReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatReportReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewReportLogic(r.Context(), svcCtx)
		resp, err := l.Report(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
