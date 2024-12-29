package group

import (
	"context"
	"im-chat/easy-chat/apps/im/rpc/imclient"
	"im-chat/easy-chat/apps/social/rpc/socialclient"
	"im-chat/easy-chat/pkg/constants"
	"im-chat/easy-chat/pkg/ctxdata"

	"im-chat/easy-chat/apps/social/api/internal/svc"
	"im-chat/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutIn 处理用户加入群组的请求（被邀请直接加入，自行申请则需要由 GroupPutInHandle 处理）
//
// 功能描述:
//   - 该方法首先调用外部服务将用户加入群组。
//   - 如果加入成功，则创建一个新的群组会话。
//   - 如果操作失败，则返回错误信息。
//
// 参数:
//   - req: `*types.GroupPutInRep` 类型，包含用户请求加入的群组ID、请求消息、请求时间和加入来源等信息。
//
// 返回值:
//   - `*types.GroupPutInResp`: 处理结果的响应。
//   - `error`: 处理过程中发生的错误。
func (l *GroupPutInLogic) GroupPutIn(req *types.GroupPutInRep) (resp *types.GroupPutInResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUId(l.ctx)

	res, err := l.svcCtx.Social.GroupPutin(l.ctx, &socialclient.GroupPutinReq{
		GroupId:    req.GroupId,
		ReqId:      uid,
		ReqMsg:     req.ReqMsg,
		ReqTime:    req.ReqTime,
		JoinSource: int32(req.JoinSource),
	})

	if res.GroupId == "" {
		return nil, err
	}

	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uid,
		RecvId:   res.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return nil, err
}
