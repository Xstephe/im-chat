package user

import (
	"im-chat/easy-chat/apps/im/ws/internal/svc"
	"im-chat/easy-chat/apps/im/ws/websocket"
)

// OnLine 处理 WebSocket 消息，向客户端发送在线用户列表。
//
// 该函数返回一个 websocket.HandlerFunc 处理函数，用于接收并处理请求在线用户列表的消息。
// 它从 WebSocket 服务器中获取所有在线用户的列表，并将该列表发送到请求的客户端。
// 如果消息发送过程中出现错误，将记录错误信息。
//
// 参数:
//   - svc: 包含服务上下文的 *svc.ServiceContext，用于访问服务相关功能。
//
// 返回:
//   - websocket.HandlerFunc: 处理 WebSocket 消息的处理函数。
func OnLine(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		uids := srv.GetUsers() //获取所有
		u := srv.GetUsers(conn)
		err := srv.Send(websocket.NewMessage(u[0], uids), conn)
		srv.Info("err ", err)
	}
}
