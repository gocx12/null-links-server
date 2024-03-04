package common

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"null-links/http_service/internal/logic/common"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadFileReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := common.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
