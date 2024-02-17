package user

import (
	"net/http"

	"null-links/http_service/internal/logic/user"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetValidationCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetValidationCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Debug("error params, error: ", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewGetValidationCodeLogic(r.Context(), svcCtx)
		resp, err := l.GetValidationCode(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
