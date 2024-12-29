package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"im-chat/easy-chat/apps/im/immodels"
	"im-chat/easy-chat/apps/im/ws/ws"
	"im-chat/easy-chat/apps/task/mq/internal/svc"
	"im-chat/easy-chat/apps/task/mq/mq"
	"im-chat/easy-chat/pkg/bitmap"
)

// MsgChatTransfer 处理聊天消息的转发。
//
// 该结构体嵌套了 baseMsgTransfer，并用于处理从消息队列中消费的聊天消息。
type MsgChatTransfer struct {
	*baseMsgTransfer
}

// NewMsgChatTransfer 创建一个新的 MsgChatTransfer 实例。
//
// 该函数用于初始化并返回一个新的消息聊天转发器实例，它封装了基本消息转发的功能。
//
// 参数:
//   - svc: 服务上下文对象，包含了服务配置和依赖的服务实例。
//
// 返回值:
//   - *MsgChatTransfer: 初始化好的消息聊天转发器实例。
func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewMsgTransfer(svc),
	}
}

// Consume 处理从消息队列中消费的聊天消息。
//
// 该方法从消息队列中获取的数据进行反序列化、记录日志，并将消息转发给目标用户。
//
// 参数:
//   - key: 消息队列中的键值，通常用于标识消息。
//   - value: 消息队列中的值，包含聊天消息的详细数据，以 JSON 格式存储。
//
// 返回值:
//   - error: 如果在处理过程中出现错误，返回相应的错误；否则返回 nil。
func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key : ", key, " value : ", value)

	var (
		data mq.MsgChatTransfer
		//ctx  = context.Background()
		msgID = primitive.NewObjectID()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := m.addChatLog(ctx, msgID, &data); err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		MsgId:          data.MsgId,
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatLog := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecords.Export()

	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
