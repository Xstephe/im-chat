package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-chat/easy-chat/apps/user/api/internal/logic/user"
	"im-chat/easy-chat/apps/user/api/internal/svc"
	"im-chat/easy-chat/apps/user/api/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
