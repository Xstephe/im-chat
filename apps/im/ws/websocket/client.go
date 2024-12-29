package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
)

// Client 表示 WebSocket 客户端，在kafka中消费。
//
// 该接口定义了 WebSocket 客户端应实现的方法，包括关闭连接、发送消息和读取消息。
type Client interface {
	Close() error

	Send(v any) error
	Read(v any) error
}

type client struct {
	*websocket.Conn
	host string

	opt dailOption
}

// NewClient 创建一个新的 WebSocket 客户端。
//
// 该函数用于创建一个新的 WebSocket 客户端实例，初始化连接到指定的 WebSocket 服务器。
// 如果连接失败，会导致程序 panic。
//
// 参数:
//   - host: WebSocket 服务器的主机地址。
//   - opts: 可选的拨号选项，用于配置 WebSocket 连接。
//
// 返回:
//   - *client: 新创建的 WebSocket 客户端实例。
func NewClient(host string, opts ...DailOptions) *client {
	opt := newDailOptions(opts...)

	c := client{
		Conn: nil,
		host: host,
		opt:  opt,
	}

	conn, err := c.dail()
	if err != nil {
		panic(err)
	}

	c.Conn = conn
	return &c
}

// dial 与 WebSocket 服务器建立连接。
//
// 该方法用于与 WebSocket 服务器建立连接，并返回一个 WebSocket 连接实例。
// 如果连接失败，则返回错误。
//
// 返回:
//   - *websocket.Conn: 成功建立的 WebSocket 连接。
//   - error: 连接过程中发生的错误（如果有的话）。
func (c *client) dail() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.pattern}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return conn, err
}

// Send 序列化并发送消息到 WebSocket。
//
// 该方法将消息对象序列化为 JSON 格式，并通过 WebSocket 连接发送。
// 如果发送失败，会尝试重新连接并重新发送消息。
//
// 参数:
//   - v: 要发送的消息对象，可以是任意类型。
//
// 返回:
//   - error: 发送消息过程中发生的错误（如果有的话）。
func (c *client) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = c.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		return nil
	}
	// todo: 再增加一个重连发送
	conn, err := c.dail()
	if err != nil {
		return err
	}
	c.Conn = conn
	return c.WriteMessage(websocket.TextMessage, data)
}

// Read 从 WebSocket 读取消息并反序列化。
//
// 该方法从 WebSocket 连接中读取消息，并将其反序列化为指定的对象类型。
// 如果读取或反序列化过程中发生错误，则返回错误。
//
// 参数:
//   - v: 用于接收反序列化后的消息对象，可以是任意类型。
//
// 返回:
//   - error: 读取消息过程中发生的错误（如果有的话）
func (c *client) Read(v any) error {
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}

	return json.Unmarshal(msg, v)
}
