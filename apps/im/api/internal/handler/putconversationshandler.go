package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-chat/easy-chat/apps/im/api/internal/logic"
	"im-chat/easy-chat/apps/im/api/internal/svc"
	"im-chat/easy-chat/apps/im/api/internal/types"
)

func putConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewPutConversationsLogic(r.Context(), svcCtx)
		resp, err := l.PutConversations(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
