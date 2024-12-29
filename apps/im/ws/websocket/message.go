package websocket

import "time"

type FrameType uint8

// 消息类型，采用grpc的写法
const (
	FrameData  FrameType = 0x0
	FramePing  FrameType = 0x1
	FrameAck   FrameType = 0x2
	FrameNoAck FrameType = 0x3
	FrameErr   FrameType = 0x9

	//FrameHeaders      FrameType = 0x1
	//FramePriority     FrameType = 0x2
	//FrameRSTStream    FrameType = 0x3
	//FrameSettings     FrameType = 0x4
	//FramePushPromise  FrameType = 0x5
	//FrameGoAway       FrameType = 0x7
	//FrameWindowUpdate FrameType = 0x8
	//FrameContinuation FrameType = 0x9
)

// msg , id, seq
type Message struct {
	FrameType `json:"frameType"`
	Id        string      `json:"id"`
	AckSeq    int         `json:"ackSeq"`
	ackTime   time.Time   `json:"ackTime"`
	errCount  int         `json:"errCount"`
	Method    string      `json:"method"`
	FormId    string      `json:"formId"`
	Data      interface{} `json:"data"` // map[string]interface{}
}

// NewMessage 创建一个新的数据消息。
//
// 该函数用于创建一个包含数据的消息对象。消息的类型被设置为 `FrameData`，
// `FormId` 是消息的发起者标识，`Data` 是消息的实际内容。
//
// 参数:
//   - formId: 消息发起者的标识符，用于标识发送消息的来源。
//   - data: 消息的实际内容，可以是任何类型的值。
//
// 返回值:
//   - *Message: 返回创建好的数据消息对象。
func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    formId,
		Data:      data,
	}
}

// NewErrMessage 创建一个新的错误消息。
//
// 该函数用于创建一个包含错误信息的消息对象。消息的类型被设置为 `FrameErr`，
// `Data` 字段包含了错误信息的字符串表示。
//
// 参数:
//   - err: 错误对象，将其错误信息转换为字符串用于消息内容。
//
// 返回值:
//   - *Message: 返回创建好的错误消息对象。
func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
