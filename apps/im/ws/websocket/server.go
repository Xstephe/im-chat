package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"time"

	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type AckType int

const (
	NoAck AckType = iota
	OnlyAck
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}

	return "NoAck"
}

// Server 表示 WebSocket 服务器的实现。
//
// 字段:
//   - routes: map[string]HandlerFunc
//     存储与请求方法对应的处理函数的路由表，每个请求方法都映射到一个特定的 `HandlerFunc`。
//   - addr: string
//     服务器监听的地址，表示 WebSocket 服务器将在哪个地址和端口上监听连接。
//   - patten: string
//     WebSocket 连接的路径模式，用于匹配客户端连接的路径。
//   - opt: *websocketOption
//     WebSocket 连接的配置选项，定义了各种 WebSocket 相关的配置参数。
//   - upgrader: websocket.Upgrader
//     WebSocket 协议升级器，用于将 HTTP 连接升级为 WebSocket 连接。
//   - Logger: logx.Logger
//     日志记录器，用于记录服务器的日志信息，包括错误、信息和调试日志。
//   - connToUser: map[*Conn]string
//     连接到用户映射表，将每个 WebSocket 连接映射到其对应的用户 ID。
//   - userToConn: map[string]*Conn
//     用户到连接映射表，将每个用户 ID 映射到其当前的 WebSocket 连接。
//   - TaskRunner: *threading.TaskRunner
//     任务运行器，用于管理和执行异步任务。
//   - RWMutex: sync.RWMutex
//     读写互斥锁，用于保护连接和用户映射表的并发读写操作。
//   - authentication: Authentication
//     鉴权接口，负责处理 WebSocket 连接的鉴权逻辑。
type Server struct {
	sync.RWMutex

	*threading.TaskRunner

	opt            *serverOption
	authentication Authentication

	routes map[string]HandlerFunc
	addr   string
	patten string

	connToUser map[*Conn]string
	userToConn map[string]*Conn

	upgrader websocket.Upgrader
	logx.Logger
}

// NewServer 创建一个新的服务器实例
//
// 该函数用于创建一个新的 `Server` 实例。它接收一个地址和可选的服务器配置选项，
// 然后返回一个初始化的 `Server` 实例。服务器实例配置了WebSocket处理、日志记录、连接
// 管理等功能。
//
// 参数:
//   - addr: 服务器监听的地址。
//   - opts: 可选的服务器配置选项，用于设置服务器的行为和特性。
//
// 返回:
//   - *Server: 初始化后的 `Server` 实例。
func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	return &Server{
		routes: make(map[string]HandlerFunc),
		addr:   addr,
		patten: opt.patten,
		opt:    &opt,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},

		authentication: opt.Authentication,

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		Logger:     logx.WithContext(context.Background()),
		TaskRunner: threading.NewTaskRunner(opt.concurrency),
	}
}

// ServerWs 处理 WebSocket 连接请求。
//
// 该方法处理 WebSocket 连接的建立。首先创建一个 WebSocket 连接对象，进行鉴权，如果鉴权失败则发送错误消息并关闭连接。
// 如果鉴权通过，将连接记录到服务器，并启动处理该连接的任务。
//
// 参数:
//   - w: HTTP 响应写入器，用于向客户端发送数据。
//   - r: HTTP 请求对象，包含请求信息。
func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	conn := NewConn(s, w, r)
	if conn != nil { // todo 这里为什么是等于，感觉是不等于
		return
	}
	//conn, err := s.upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	s.Errorf("upgrade err %v", err)
	//	return
	//}

	//对连接的鉴权
	if !s.authentication.Auth(w, r) {
		//conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("不具备访问权限")))
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprint("不具备访问权限")}, conn)
		conn.Close()
		return
	}

	// 记录连接
	s.addConn(conn, r)

	// 处理连接
	go s.handlerConn(conn)
}

// handlerConn 根据连接对象进行任务处理。
//
// 该方法处理 WebSocket 连接的消息任务。它从连接中读取消息，解析消息，并根据配置处理消息的 ACK 确认。
// 如果启用了 ACK 机制且消息需要确认，方法会将消息放入待确认的消息队列；否则，直接将消息传递给处理函数。
//
// 参数:
//   - conn: WebSocket 连接对象，用于接收和发送消息。
func (s *Server) handlerConn(conn *Conn) {

	uids := s.GetUsers(conn)
	conn.Uid = uids[0]

	// 处理任务
	go s.handlerWrite(conn)

	if s.isAck(nil) {
		go s.readAck(conn)
	}

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			//关闭连接
			s.Close(conn)
			return
		}
		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(msg))
			s.Close(conn)
			return
		}

		// 依据消息进行处理
		if s.isAck(&message) {
			s.Infof("conn message read ack msg %v", message)
			conn.appendMsgMq(&message)
		} else {
			conn.message <- &message
		}
	}
}

// isAck 判断当前消息是否需要进行 ACK 确认。
//
// 该方法根据消息的状态和服务器配置判断是否需要对消息进行 ACK 确认。
// 如果消息为空且服务器配置要求 ACK，则返回 true；
// 如果消息不为空且其 FrameType 为 FrameNoAck，则返回 false；
// 否则，根据服务器的 ACK 配置和消息的 FrameType 返回是否需要 ACK 确认。
//
// 参数:
//   - message: 要判断的消息对象，如果为 nil，则根据服务器的 ACK 配置判断是否需要确认。
//
// 返回:
//   - bool: 是否需要进行 ACK 确认。
func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck
}

// 读取消息的ack
func (s *Server) readAck(conn *Conn) {
	for {
		select {
		case <-conn.done:
			s.Infof("close message ack uid %v ", conn.Uid)
			return
		default:
		}

		// 从队列中读取新的消息
		conn.messageMu.Lock()
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			// 增加睡眠
			time.Sleep(100 * time.Microsecond)
			continue
		}

		// 读取第一条
		message := conn.readMessage[0]

		// 判断ack的方式
		switch s.opt.ack {
		case OnlyAck:
			// 直接给客户端回复
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn)
			// 进行业务处理
			// 把消息从队列中移除
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()

			conn.message <- message
		case RigorAck:
			// 先回
			if message.AckSeq == 0 {
				// 还未确认
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].ackTime = time.Now()
				s.Send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn)
				s.Infof("message ack RigorAck send mid %v, seq %v , time%v", message.Id, message.AckSeq,
					message.ackTime)
				conn.messageMu.Unlock()
				continue
			}

			// 再验证

			// 1. 客户端返回结果，再一次确认
			// 得到客户端的序号
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				// 确认
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack RigorAck success mid %v", message.Id)
				continue
			}

			// 2. 客户端没有确认，考虑是否超过了ack的确认时间
			val := s.opt.ackTimeout - time.Since(message.ackTime)
			if !message.ackTime.IsZero() && val <= 0 {
				//		2.2 超过结束确认
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}
			//		2.1 未超过，重新发送
			conn.messageMu.Unlock()
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn)
			// 睡眠一定的时间
			time.Sleep(3 * time.Second)
		}
	}
}

// handleWrite 处理并分发消息任务。
//
// 该方法用于处理连接中的消息，并根据消息的 FrameType 分发到相应的处理器。
// 它在循环中接收消息并根据消息的类型进行处理。
// 如果消息需要 ACK 确认，则清除消息确认状态。
// 处理过程包括以下几种情况：
// - 对于 Ping 消息，直接发送 Ping 响应。
// - 对于 Data 消息，根据消息的 Method 执行对应的处理器。
// - 处理完成后，如果消息需要 ACK 确认，则从连接的消息队列中删除该消息的确认状态。
//
// 参数:
//   - conn: 连接对象，消息将从该连接的消息队列中读取。
func (s *Server) handlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 连接关闭
			return
		case message := <-conn.message:
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{FrameType: FramePing}, conn)
			case FrameData:
				// 根据请求的method分发路由并执行
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)}, conn)
					//conn.WriteMessage(&Message{}, []byte(fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)))
				}
			}

			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}

// addConn 存储 WebSocket 连接并与用户 ID 关联。
//
// 该方法用于将新的 WebSocket 连接添加到服务器中，并将其与用户 ID 进行关联。
// 如果用户 ID 对应的连接已存在，则关闭之前的连接。
// 将新的连接与用户 ID 关联后，更新服务器的连接映射。
//
// 参数:
//   - conn: 要添加的 WebSocket 连接。
//   - req: HTTP 请求，用于获取用户 ID。
func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否之前登入过
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

// GetConn 根据用户 ID 获取 WebSocket 连接。
//
// 该方法用于根据用户 ID 从服务器的用户到连接的映射中获取对应的 WebSocket 连接。
// 如果找不到对应的连接，返回 nil。
//
// 参数:
//   - uid: 用户的 ID。
//
// 返回:
//   - *Conn: 对应用户 ID 的 WebSocket 连接；如果未找到，则返回 nil。
func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

// GetConns 根据用户 ID 列表获取 WebSocket 连接列表。
//
// 该方法用于根据用户 ID 列表从服务器的用户到连接的映射中获取对应的 WebSocket 连接列表。
// 如果用户 ID 列表为空，返回一个空的连接列表。
//
// 参数:
//   - uids: 用户的 ID 列表。
//
// 返回:
//   - []*Conn: 对应用户 ID 的 WebSocket 连接列表；如果用户 ID 列表为空，则返回 nil。
func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

// GetUsers 获取与指定连接关联的用户 ID 列表。
//
// 该方法用于从服务器获取与连接相关的用户 ID 列表。
// 如果没有指定连接，则返回所有用户的 ID 列表；如果指定了连接，则返回这些连接对应的用户 ID 列表。
//
// 参数:
//   - conns: 要获取用户 ID 的 WebSocket 连接列表。如果为空，则返回所有用户的 ID 列表。
//
// 返回:
//   - []string: 用户 ID 列表。如果没有连接，则返回所有用户的 ID 列表。
func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

// Close 关闭指定的 WebSocket 连接并从服务器中移除。
//
// 该方法用于关闭 WebSocket 连接并从服务器的连接映射中删除该连接。
// 如果连接已经被关闭（即用户 ID 为空），则不执行任何操作。
// 关闭连接后，将从 `connToUser` 和 `userToConn` 映射中移除相关条目。
//
// 参数:
//   - conn: 要关闭的 WebSocket 连接。
func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	conn.Close()
}

// SendByUserId 向指定的用户 ID 发送消息。
//
// 该方法用于将消息发送到指定的用户连接。首先，根据传入的用户 ID 获取对应的连接，然后将消息发送到这些连接中。
// 如果没有指定用户 ID，则不执行任何操作。
//
// 参数:
//   - msg: 要发送的消息，可以是任何类型的数据。
//   - sendIds: 目标用户的 ID 列表，消息将发送到这些用户的连接中。
//
// 返回:
//   - error: 如果发送过程中发生错误，返回错误信息；如果成功，则返回 nil。
func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

// Send 向指定的连接发送消息。
//
// 该方法用于将消息发送到一个或多个 WebSocket 连接。
// 首先将消息序列化为 JSON 格式，然后遍历连接列表，将消息通过 WebSocket 发送到每个连接中。
// 如果没有指定连接，则不执行任何操作。
// 如果在发送过程中发生错误，方法将立即返回该错误；如果成功，则返回 nil。
//
// 参数:
//   - msg: 要发送的消息，可以是任何类型的数据。
//   - conns: 要发送消息的连接列表。
//
// 返回:
//   - error: 如果发送过程中发生错误，返回错误信息；如果成功，则返回 nil。
func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

// AddRoutes 注册 WebSocket 路由处理函数。
//
// 该方法将一组路由注册到服务器中，用于处理 WebSocket 请求消息。
// 每个路由包括一个方法名和一个处理函数，根据请求消息的方法名将请求分发到对应的处理函数。
//
// 参数:
//   - rs: 包含路由信息的列表，每个路由包含方法名和处理函数。
func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

// Start 启动服务器
//
// 该方法用于启动HTTP服务器并开始监听指定的地址。它将处理所有传入的请求，并调用
// `ServerWs` 方法处理WebSocket连接。启动后，服务器将会持续运行，直到出现错误或
// 被手动停止。
func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

// Stop 停止服务器
//
// 该方法用于停止正在运行的服务器。它会打印一条停止服务的消息。请注意，该方法
// 只是打印了停止服务的消息，实际停止服务的操作可能需要额外的实现。
func (s *Server) Stop() {
	fmt.Println("停止服务")
}
