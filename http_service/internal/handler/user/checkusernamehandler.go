package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"null-links/http_service/internal/logic/user"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
)

func CheckUsernameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckUsernameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewCheckUsernameLogic(r.Context(), svcCtx)
		resp, err := l.CheckUsername(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
