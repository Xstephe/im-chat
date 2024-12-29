package push

import (
	"github.com/mitchellh/mapstructure"
	"im-chat/easy-chat/apps/im/ws/internal/svc"
	"im-chat/easy-chat/apps/im/ws/websocket"
	"im-chat/easy-chat/apps/im/ws/ws"
	"im-chat/easy-chat/pkg/constants"
)

// Push 处理 WebSocket 消息，转发推送消息，由 kafka 消息队列远程调用。
//
// 该函数返回一个 websocket.HandlerFunc 处理函数，用于接收并处理推送消息。
// 它将 WebSocket 消息解码为 ws.Push 结构体，并根据聊天类型将消息推送到目标用户。
// 如果消息解码失败，或推送过程中出现错误，将通过 WebSocket 向客户端发送错误信息。
//
// 参数:
//   - svc: 包含服务上下文的 *svc.ServiceContext，用于访问服务相关功能。
//
// 返回:
//   - websocket.HandlerFunc: 处理 WebSocket 消息的处理函数。
func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err))
			return
		}

		switch data.ChatType {
		case constants.SingleChatType:
			single(srv, &data, data.RecvId)
		case constants.GroupChatType:
			group(srv, &data)
		}
	}
}

// single 处理单聊消息的推送。
//
// 该函数根据接收者ID从服务器获取连接，并将消息推送给接收者。
// 如果目标用户离线，当前实现没有处理离线用户的逻辑。
// 如果推送过程中出现错误，记录错误日志。
//
// 参数:
//   - srv: WebSocket 服务器实例。
//   - data: 包含推送消息的数据结构体。
//   - recvId: 接收者用户ID。
//
// 返回:
//   - error: 发生的错误（如果有的话），返回nil表示推送成功。
func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	// 发送的目标
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// todo: 目标离线
		return nil
	}

	srv.Infof("push msg %v", data)

	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			ReadRecords: data.ReadRecords,
			MsgId:       data.MsgId,
			MType:       data.MType,
			Content:     data.Content,
		},
	}), rconn)
}

func group(srv *websocket.Server, data *ws.Push) error {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, data.RecvId)
			})
		}(id)
	}
	return nil
}
