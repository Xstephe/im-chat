package friend

import (
	"context"
	"im-chat/easy-chat/apps/social/rpc/socialclient"
	"im-chat/easy-chat/pkg/ctxdata"

	"im-chat/easy-chat/apps/social/api/internal/svc"
	"im-chat/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendPutIn 处理用户添加好友的请求
//
// 功能描述:
//   - 从上下文中获取当前用户ID
//   - 调用服务层接口发起好友添加请求，传递请求用户ID、被请求用户ID、请求消息和请求时间
//
// 参数:
//   - req: `*types.FriendPutInReq` 类型，包含添加好友请求的相关信息
//   - `UserId`: 当前用户ID，即请求添加好友的用户
//   - `ReqUid`: 被请求添加为好友的用户ID
//   - `ReqMsg`: 附带的请求消息
//   - `ReqTime`: 请求时间
//
// 返回值:
//   - `*types.FriendPutInResp`: 响应对象，包含操作结果。
//   - `error`: 如果在处理请求过程中发生错误，则返回相应的错误信息。
func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUId(l.ctx)

	_, err = l.svcCtx.Social.FriendPutIn(l.ctx, &socialclient.FriendPutInReq{
		UserId:  uid,
		ReqUid:  req.UserId,
		ReqMsg:  req.ReqMsg,
		ReqTime: req.ReqTime,
	})
	if err != nil {
		return nil, err
	}

	return &types.FriendPutInResp{}, nil
}
