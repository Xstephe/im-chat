package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-chat/easy-chat/apps/im/api/internal/logic"
	"im-chat/easy-chat/apps/im/api/internal/svc"
	"im-chat/easy-chat/apps/im/api/internal/types"
)

func setUpUserConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetUpUserConversationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSetUpUserConversationLogic(r.Context(), svcCtx)
		resp, err := l.SetUpUserConversation(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
