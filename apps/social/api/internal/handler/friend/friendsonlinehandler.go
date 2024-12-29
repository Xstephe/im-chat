package friend

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-chat/easy-chat/apps/social/api/internal/logic/friend"
	"im-chat/easy-chat/apps/social/api/internal/svc"
	"im-chat/easy-chat/apps/social/api/internal/types"
)

func FriendsOnlineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendsOnlineReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := friend.NewFriendsOnlineLogic(r.Context(), svcCtx)
		resp, err := l.FriendsOnline(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
