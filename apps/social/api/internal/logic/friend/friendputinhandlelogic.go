package friend

import (
	"context"
	"im-chat/easy-chat/apps/social/rpc/socialclient"
	"im-chat/easy-chat/pkg/ctxdata"

	"im-chat/easy-chat/apps/social/api/internal/svc"
	"im-chat/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendPutInHandle 处理好友申请
//
// 功能描述:
//   - 从上下文中获取当前用户ID
//   - 调用服务层接口处理好友申请，包括审批或拒绝操作
//   - 返回操作结果
//
// 参数:
//   - req: `*types.FriendPutInHandleReq` 类型，包含处理好友申请的请求信息
//   - `FriendReqId`: 好友申请记录的唯一标识符
//   - `HandleResult`: 处理结果，表示审批通过或拒绝
//
// 返回值:
//   - `*types.FriendPutInHandleResp`: 响应对象，表示处理好友申请的结果
//   - `error`: 如果在处理好友申请过程中发生错误，则返回相应的错误信息
func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (resp *types.FriendPutInHandleResp, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.Social.FriendPutInHandle(l.ctx, &socialclient.FriendPutInHandleReq{
		FriendReqId:  req.FriendReqId,
		UserId:       ctxdata.GetUId(l.ctx),
		HandleResult: req.HandleResult,
	})
	return
}
