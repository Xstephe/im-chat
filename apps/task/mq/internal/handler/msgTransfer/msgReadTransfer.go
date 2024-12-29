package msgTransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"im-chat/easy-chat/apps/im/ws/ws"
	"im-chat/easy-chat/apps/task/mq/internal/svc"
	"im-chat/easy-chat/apps/task/mq/mq"
	"im-chat/easy-chat/pkg/bitmap"
	"im-chat/easy-chat/pkg/constants"
	"sync"
	"time"
)

var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

type MsgReadTransfer struct {
	*baseMsgTransfer

	cache.Cache

	mu sync.RWMutex

	groupMsgs map[string]*groupMsgRead
	push      chan *ws.Push
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {

	m := &MsgReadTransfer{
		baseMsgTransfer: NewMsgTransfer(svc),
		groupMsgs:       make(map[string]*groupMsgRead, 1),
		push:            make(chan *ws.Push, 1),
	}
	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}

		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime) * time.Second
		}
	}

	go m.transfer()
	return m
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	m.Info("MsgReadTransfer", value)

	var (
		data mq.MsgMarkRead
		//ctx  = context.Background()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	//业务处理 -- 更新
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMarkRead,
		ReadRecords:    readRecords,
	}
	switch data.ChatType {
	case constants.SingleChatType:
		//直接推送
		m.push <- push
	case constants.GroupChatType:
		//判断是否开启合并消息的处理
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.push <- push
		}
		m.mu.Lock()
		defer m.mu.Unlock()

		push.SendId = "" // 群聊不需要发送者ID，统一清空
		// 已经有消息
		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			// 和并请求
			m.Infof("merge push: %v", push.ConversationId)
			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			// 没有记录，创建
			m.Infof("create merge push %v", push.ConversationId)
			m.groupMsgs[push.ConversationId] = newGroupMsgRead(push, m.push)
		}
	}
	return nil

}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	res := make(map[string]string)

	chatLogs, err := m.svcCtx.ChatLogModel.ListByIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}

	//处理已读
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}
		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

		err = m.svcCtx.ChatLogModel.UpdateMakeRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}
	return res, nil

}

// 异步处理消息发送
func (m *MsgReadTransfer) transfer() {
	for push := range m.push {
		if push.RecvId != "" || len(push.RecvIds) > 0 {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("transfer err: %s", err.Error())
			}
		}
		if push.ChatType == constants.SingleChatType {
			continue
		}
		// 不采用合并推送
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}
		// 清空数据
		m.mu.Lock()
		if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
			m.groupMsgs[push.ConversationId].Clear()
			delete(m.groupMsgs, push.ConversationId)
		}
		m.mu.Unlock()
	}
}
