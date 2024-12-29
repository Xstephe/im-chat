package svc

import (
	"im-chat/easy-chat/apps/im/immodels"
	"im-chat/easy-chat/apps/im/ws/internal/config"
	"im-chat/easy-chat/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config

	immodels.ChatLogModel
	mqclient.MsgChatTransferClient
	mqclient.MsgReadTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		MsgReadTransferClient: mqclient.NewmsgReadTransferClient(c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic),
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
