package conversation

import (
	"github.com/mitchellh/mapstructure"
	"im-chat/easy-chat/apps/im/ws/internal/svc"
	"im-chat/easy-chat/apps/im/ws/websocket"
	"im-chat/easy-chat/apps/im/ws/ws"
	"im-chat/easy-chat/apps/task/mq/mq"
	"im-chat/easy-chat/pkg/constants"
	"im-chat/easy-chat/pkg/wuid"
	"time"
)

// Chat 处理 WebSocket 消息，进行聊天消息的转发。
//
// 该函数返回一个 websocket.HandlerFunc 处理函数，用于接收并处理聊天消息。
// 它将 WebSocket 消息解码为 ws.Chat 结构体，若消息未指定会话ID，则根据聊天类型生成会话ID。
// 处理完成后，将聊天消息推送到消息聊天传输客户端进行处理。
// 如果解码或消息处理失败，将通过 WebSocket 向客户端发送错误信息。
//
// 参数:
//   - svc: 包含服务上下文的 *svc.ServiceContext，用于访问消息聊天传输客户端。
//
// 返回:
//   - websocket.HandlerFunc: 处理 WebSocket 消息的处理函数。
func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 私聊
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.SingleChatType:
				data.ConversationId = wuid.CombineId(conn.Uid, data.RecvId)
			case constants.GroupChatType:
				data.ConversationId = data.RecvId
			}

		}

		err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			SendTime:       time.Now().UnixMilli(),
			MType:          data.Msg.MType,
			Content:        data.Msg.Content,
			MsgId:          msg.Id,
		})
		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		//err := logic.NewConversation(context.Background(), srv, svc).SingleChat(&data, conn.Uid)
		//if err != nil {
		//	srv.Send(websocket.NewErrMessage(err), conn)
		//	return
		//}
		//srv.SendByUserId(websocket.NewMessage(conn.Uid, ws.Chat{
		//	ConversationId: data.ConversationId,
		//	ChatType:       data.ChatType,
		//	SendId:         conn.Uid,
		//	RecvId:         data.RecvId,
		//	SendTime:       time.Now().UnixMilli(),
		//	Msg:            data.Msg,
		//}), data.RecvId)

	}
}

func MarkRead(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// todo: 已读未读处理
		var data ws.MarkRead
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		err := svc.MsgReadTransferClient.Push(&mq.MsgMarkRead{
			ChatType:       data.ChatType,
			RecvId:         data.RecvId,
			SendId:         conn.Uid,
			ConversationId: data.ConversationId,
			MsgIds:         data.MsgIds,
		})
		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

	}
}
