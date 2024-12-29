package msgTransfer

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"im-chat/easy-chat/apps/im/ws/websocket"
	"im-chat/easy-chat/apps/im/ws/ws"
	"im-chat/easy-chat/apps/social/rpc/socialclient"
	"im-chat/easy-chat/apps/task/mq/internal/svc"
	"im-chat/easy-chat/pkg/constants"
)

type baseMsgTransfer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMsgTransfer(svcCtx *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(context.Background()),
	}
}

// 用于转发消息
func (m *baseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.SingleChatType:
		err = m.single(ctx, data)
	case constants.GroupChatType:
		err = m.group(ctx, data)
	}
	return err
}

// 私聊
func (m *baseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

// 群聊
func (m *baseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	//就要查询群的用户
	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}

	data.RecvIds = make([]string, 0, len(users.List))

	for _, members := range users.List {
		if members.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, members.UserId)
	}

	// 推送消息
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}
