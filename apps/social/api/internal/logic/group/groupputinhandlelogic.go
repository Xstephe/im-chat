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

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutInHandle 处理群组加入请求
//
// 功能描述:
//   - 处理群组加入请求，根据请求的处理结果更新请求状态
//   - 如果请求被批准，则建立一个新的用户和群组的聊天会话
//
// 参数:
//   - req: `*types.GroupPutInHandleReq` 类型，包含处理群组请求所需的信息，包括群组请求ID、群组ID、处理结果等
//
// 返回值:
//   - `*types.GroupPutInHandleResp`: 空响应对象
//   - `error`: 如果在处理过程中发生错误，则返回相应的错误信息
func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleRep) (resp *types.GroupPutInHandleResp, err error) {
	// todo: add your logic here and delete this line

	uid := ctxdata.GetUId(l.ctx)
	res, err := l.svcCtx.Social.GroupPutInHandle(l.ctx, &socialclient.GroupPutInHandleReq{
		GroupReqId:   req.GroupReqId,
		GroupId:      req.GroupId,
		HandleUid:    uid,
		HandleResult: req.HandleResult,
	})

	if constants.HandlerResult(req.HandleResult) != constants.PassHandlerResult {
		return
	}

	// todo: 通过后的业务

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
