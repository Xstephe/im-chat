package handler

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"im-chat/easy-chat/apps/task/mq/internal/handler/msgTransfer"
	"im-chat/easy-chat/apps/task/mq/internal/svc"
)

type Listen struct {
	svc *svc.ServiceContext
}

func NewListen(svc *svc.ServiceContext) *Listen {
	return &Listen{svc: svc}
}

func (l *Listen) Services() []service.Service {
	return []service.Service{
		// todo: 此处可以加载多个消费者
		kq.MustNewQueue(l.svc.Config.MsgReadTransfer, msgTransfer.NewMsgReadTransfer(l.svc)),
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer, msgTransfer.NewMsgChatTransfer(l.svc)),
	}
}
