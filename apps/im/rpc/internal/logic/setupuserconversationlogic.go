package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"im-chat/easy-chat/apps/im/immodels"
	"im-chat/easy-chat/apps/im/rpc/im"
	"im-chat/easy-chat/apps/im/rpc/internal/svc"
	"im-chat/easy-chat/pkg/constants"
	"im-chat/easy-chat/pkg/wuid"
	"im-chat/easy-chat/pkg/xerr"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetUpUserConversation 建立会话: 群聊, 私聊
//
// 该方法负责创建或设置用户会话，支持单聊和群聊两种类型。
// 根据聊天类型，创建或更新相应的会话记录。
//
// 参数:
// - in: 包含建立会话所需的请求参数，主要包括发送者ID、接收者ID和聊天类型。
//
// 返回:
// - *im.SetUpUserConversationResp: 返回会话设置的响应结构体。
// - error: 发生的错误（如果有的话），返回nil表示操作成功。
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// todo: add your logic here and delete this line

	var res im.SetUpUserConversationResp
	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		// 生成会话的id
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		// 验证是否建立过会话
		conversationRes, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			// 建立会话
			if err == immodels.ErrNotFound {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.SingleChatType,
				})

				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v", err)
				}
			} else {
				return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, conversationId)
			}
		} else if conversationRes != nil {
			return &res, nil
		}
		// 建立两者的会话
		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, true)
		if err != nil {
			return nil, err
		}
		err = l.setUpUserConversation(conversationId, in.RecvId, in.SendId, constants.SingleChatType, false)
		if err != nil {
			return nil, err
		}
	case constants.GroupChatType:
		err := l.setUpUserConversation(in.RecvId, in.SendId, in.RecvId, constants.GroupChatType, true)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string,
	chatType constants.ChatType, isShow bool) error {
	// 用户的会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == immodels.ErrNotFound {
			conversations = &immodels.Conversations{
				ID:               primitive.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*immodels.Conversation),
			}
		} else {
			return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, userId)
		}
	}

	// 更新会话记录,会话记录是否存在过
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}

	// 添加会话记录
	conversations.ConversationList[conversationId] = &immodels.Conversation{
		ConversationId: conversationId,
		ChatType:       constants.SingleChatType,
		IsShow:         isShow,
	}

	// 更新
	err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v, req %v", err, conversations)
	}
	return nil
}
