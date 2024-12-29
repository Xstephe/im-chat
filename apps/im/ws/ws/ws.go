package ws

import "im-chat/easy-chat/pkg/constants"

type (
	// Msg 表示一个基础消息的结构体。
	//
	// 该结构体包含消息的唯一标识符、已读记录、消息类型和消息内容
	Msg struct {
		constants.MType `mapstructure:"mType"`
		Content         string            `mapstructure:"content"`
		MsgId           string            `mapstructure:"msgId"`
		ReadRecords     map[string]string `mapstructure:"readRecords"`
	}

	// Chat 表示一个聊天消息的结构体。
	//
	// 该结构体继承了 Msg 结构体，包含了会话ID、聊天类型、发送者和接收者ID、发送时间等信息。
	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`
		Msg                `mapstructure:"msg"`
	}

	// Push 表示一个推送消息的结构体。
	//
	// 该结构体包含了推送消息所需的信息，包括会话ID、发送者和接收者ID列表、发送时间、消息内容等。
	Push struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string   `mapstructure:"sendId"`
		RecvId             string   `mapstructure:"recvId"`
		RecvIds            []string `mapstructure:"recvIds"`
		SendTime           int64    `mapstructure:"sendTime"`

		MsgId       string                `mapstructure:"msgId"`
		ReadRecords map[string]string     `mapstructure:"readRecords"`
		ContentType constants.ContentType `mapstructure:"contentType"`

		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}

	// MarkRead 表示一个标记消息已读的结构体。
	//
	// 该结构体用于处理标记消息已读的操作，包括会话ID、接收者ID和已读的消息ID列表。
	MarkRead struct {
		constants.ChatType `mapstructure:"chatType"`
		RecvId             string   `mapstructure:"recvId"`
		ConversationId     string   `mapstructure:"conversationId"`
		MsgIds             []string `mapstructure:"msgIds"`
	}
)
