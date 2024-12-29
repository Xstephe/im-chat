package logic

import (
	"context"
	"github.com/pkg/errors"
	"im-chat/easy-chat/apps/im/immodels"
	"im-chat/easy-chat/apps/im/rpc/im"
	"im-chat/easy-chat/apps/im/rpc/internal/svc"
	"im-chat/easy-chat/pkg/constants"
	"im-chat/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateGroupConversation 创建群聊会话
//
// 该方法用于创建群聊会话。如果指定的群ID已存在，则直接返回；
// 如果群ID不存在，则新建一个群聊会话，并将创建者的用户会话列表进行更新。
//
// 参数:
//   - in: 包含群聊会话创建请求的结构体，包括群ID和创建者用户ID。
//
// 返回:
//   - *im.CreateGroupConversationResp: 创建群聊会话的响应结构体，包含相关的响应数据。
//   - error: 发生的错误（如果有的话），返回nil表示操作成功。
func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	// todo: add your logic here and delete this line
	res := &im.CreateGroupConversationResp{}

	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	if err == nil {
		return res, nil
	}

	if err != immodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.FindOne err %v,req %v", err, in.GroupId)
	}

	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.Insert err %v", err)
	}

	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "NewSetUpUserConversationLogic err %v", err)
	}

	return res, nil
}
