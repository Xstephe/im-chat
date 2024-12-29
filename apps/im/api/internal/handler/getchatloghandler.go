package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-chat/easy-chat/apps/im/api/internal/logic"
	"im-chat/easy-chat/apps/im/api/internal/svc"
	"im-chat/easy-chat/apps/im/api/internal/types"
)

func getChatLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatLogReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetChatLogLogic(r.Context(), svcCtx)
		resp, err := l.GetChatLog(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
