package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"im-chat/easy-chat/apps/im/immodels"
	"im-chat/easy-chat/apps/im/rpc/im"
	"im-chat/easy-chat/apps/im/rpc/internal/svc"
	"im-chat/easy-chat/pkg/constants"
	"im-chat/easy-chat/pkg/xerr"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PutConversations 更新会话信息。
//
// 该方法更新指定用户的会话列表，将新的会话信息保存到数据库中。
// 如果用户原本有会话数据，将会合并新数据；如果没有，将会创建新的会话记录。
//
// 参数:
//   - in: 请求对象，包含需要更新的会话信息。
//
// 返回值:
//   - *im.PutConversationsResp: 响应对象，表示更新操作的结果。
//   - error: 如果在更新过程中发生错误，返回具体的错误信息；成功时返回 nil。
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	// todo: add your logic here and delete this line

	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindByUserId err %v, req %v", err, in.UserId)
	}

	if data.ConversationList == nil {
		data.ConversationList = make(map[string]*immodels.Conversation)
	}

	for s, conversation := range in.ConversationList {
		var oldTotal int
		if data.ConversationList[s] != nil {
			oldTotal = data.ConversationList[s].Total
		}

		data.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversation.ConversationId,
			ChatType:       constants.ChatType(conversation.ChatType),
			IsShow:         conversation.IsShow,
			Total:          int(conversation.Read) + oldTotal,
			Seq:            conversation.Seq,
		}
	}

	err = l.svcCtx.ConversationsModel.Update(l.ctx, data)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Update err %v, req %v", err, data)
	}

	return &im.PutConversationsResp{}, nil
}
