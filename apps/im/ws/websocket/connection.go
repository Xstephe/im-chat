package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// Conn 表示 WebSocket 连接。
//
// 该结构体定义了一个WebSocket连接的主要属性和状态，包括用户ID、WebSocket连接实例、
// 连接的管理服务器以及用于消息处理和连接状态维护的各种信息。此结构体用于管理和维护
// WebSocket连接的生命周期和消息传递。
//
// 字段:
//   - idleMu: 连接空闲状态的互斥锁，用于保护空闲时间的读写操作。
//   - Uid: 用户标识符，用于标识与该连接关联的用户。
//   - websocket.Conn: WebSocket连接实例，表示与客户端的实际WebSocket连接。
//   - s: 连接所属的WebSocket服务器，用于访问服务器相关的功能和状态。
//   - idle: 连接的空闲时间，用于检测连接的活动状态。
//   - maxConnectionIdle: 允许的最大空闲时间，超过该时间连接将被认为是超时。
//   - messageMu: 消息队列的互斥锁，用于保护消息队列的读写操作。
//   - readMessage: 读消息队列，存储尚未处理的消息。
//   - readMessageSeq: 读消息队列的序列化映射，用于按序号存储消息。
//   - message: 消息通道，用于接收和发送消息。
//   - done: 关闭连接时的信号通道，用于通知连接的结束。
type Conn struct {
	idleMu sync.Mutex

	Uid string

	*websocket.Conn
	s *Server

	idle              time.Time
	maxConnectionIdle time.Duration

	messageMu      sync.Mutex
	readMessage    []*Message
	readMessageSeq map[string]*Message

	message chan *Message

	done chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {

	var responseHeader http.Header
	if protocol := r.Header.Get("Sec-WebSocket-Protocol"); protocol != "" {
		responseHeader = http.Header{"Sec-WebSocket-Protocol": []string{protocol}}

	}

	c, err := s.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		readMessage:       make([]*Message, 0, 2),
		readMessageSeq:    make(map[string]*Message, 2),
		message:           make(chan *Message, 1),
		done:              make(chan struct{}),
	}

	go conn.keepalive()
	return conn
}

// appendMsgMq 将消息添加到消息队列中。
//
// 该方法用于将传入的消息添加到读消息队列中，并维护消息序列化映射。
// 如果消息已经存在于队列中且其确认序号小于等于之前存储的消息，则忽略该消息。
// 如果消息类型是确认消息（FrameAck），则不处理。
// 否则，将消息添加到队列并更新消息序列化映射。
//
// 参数:
//   - msg: 要添加到消息队列中的消息结构体。
func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()

	// 读队列中
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		// 已经有消息的记录，该消息已经有ack的确认
		if len(c.readMessage) == 0 {
			// 队列中没有该消息
			return
		}

		// msg.AckSeq > m.AckSeq
		if m.AckSeq >= msg.AckSeq {
			// 没有进行ack的确认, 重复
			return
		}

		c.readMessageSeq[msg.Id] = msg
		return
	}
	// 还没有进行ack的确认, 避免客户端重复发送多余的ack消息
	if msg.FrameType == FrameAck {
		return
	}

	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg

}

// ReadMessage 从 WebSocket 连接中读取消息。
//
// 该方法从WebSocket连接中读取消息，并重置连接的空闲时间。
// 如果读取消息时发生错误，则返回该错误。
// 空闲时间用于管理连接的活跃状态。
//
// 返回:
//   - messageType: 消息类型，指示消息的格式（文本或二进制）。
//   - p: 读取到的消息内容。
//   - err: 读取消息时发生的错误，如果没有错误则返回nil。
func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()

	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{}
	return
}

// WriteMessage 向 WebSocket 连接中写入消息。
//
// 该方法将消息写入WebSocket连接，并更新连接的空闲时间。
// 如果写入消息时发生错误，则返回该错误。
// 空闲时间用于管理连接的活跃状态。
//
// 参数:
//   - messageType: 消息类型，指示消息的格式（文本或二进制）。
//   - data: 要写入的消息内容。
//
// 返回:
//   - error: 写入消息时发生的错误，如果成功写入则返回nil。
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	// 方法是并不安全，加锁
	err := c.Conn.WriteMessage(messageType, data)
	c.idle = time.Now()
	return err
}

// Close 关闭 WebSocket 连接。
//
// 该方法关闭WebSocket连接并通知所有相关操作停止。
// 如果连接已经关闭，则不会重复关闭。关闭操作通过信号通道通知。
// 关闭操作返回WebSocket连接关闭的错误（如果有的话）。
//
// 返回:
//   - error: 关闭WebSocket连接时发生的错误，如果成功关闭则返回nil。
func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	return c.Conn.Close()
}

// keepalive 定期检查连接的空闲状态，确保连接在超过最大空闲时间后被优雅地关闭。
// 该方法会启动一个定时器，根据连接的最大空闲时间进行检查。
// 如果连接空闲时间超过了最大空闲时间，连接将被关闭。
// 如果连接未超过空闲时间，定时器将重置以继续监控连接状态。
// 方法会监听连接的关闭事件，以便在连接关闭时终止检查。
func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() { // The connection is non-idle.
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// The connection has been idle for a duration of keepalive.MaxConnectionIdle or more.
				// Gracefully close the connection.
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			return
		}
	}
}
