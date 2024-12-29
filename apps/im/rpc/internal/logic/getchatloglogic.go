package logic

import (
	"context"
	"github.com/pkg/errors"
	"im-chat/easy-chat/pkg/xerr"

	"im-chat/easy-chat/apps/im/rpc/im"
	"im-chat/easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetChatLog 获取会话记录。
//
// 该方法根据请求中的参数从数据库中获取聊天记录。根据是否提供了 msgId，
// 方法会选择不同的查询方式：如果 msgId 不为空，则直接查询该消息记录；
// 如果 msgId 为空，则根据时间段进行查询。查询的结果会按照时间排序，并返回符合条件的聊天记录。
//
// 参数:
//   - in: 请求对象，包含查询条件。
//
// 返回值:
//
//   - *im.GetChatLogResp: 查询结果的响应对象。
func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// todo: add your logic here and delete this line

	// 根据id
	if in.MsgId != "" {
		chatLog, err := l.svcCtx.ChatLogModel.FindOne(l.ctx, in.MsgId)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "find chatlog by msgId err %v, req %v", err, in.MsgId)
		}

		return &im.GetChatLogResp{
			List: []*im.ChatLog{
				{
					Id:             chatLog.ID.Hex(),
					ConversationId: chatLog.ConversationId,
					SendId:         chatLog.SendId,
					RecvId:         chatLog.RecvId,
					MsgType:        int32(chatLog.MsgType),
					MsgContent:     chatLog.MsgContent,
					ChatType:       int32(chatLog.ChatType),
					SendTime:       chatLog.SendTime,
				},
			},
		}, nil
	}
	// 时间段分段查询
	data, err := l.svcCtx.ChatLogModel.ListBySendTime(l.ctx, in.ConversationId, in.StartSendTime, in.EndSendTime, in.Count)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ListBySendTime err %v, req %v", err, in)
	}

	res := make([]*im.ChatLog, 0, len(data))
	for _, datum := range data {
		res = append(res, &im.ChatLog{
			Id:             datum.ID.Hex(),
			ConversationId: datum.ConversationId,
			SendId:         datum.SendId,
			RecvId:         datum.RecvId,
			MsgType:        int32(datum.MsgType),
			MsgContent:     datum.MsgContent,
			ChatType:       int32(datum.ChatType),
			SendTime:       datum.SendTime,
		})
	}

	return &im.GetChatLogResp{
		List: res,
	}, nil
}
