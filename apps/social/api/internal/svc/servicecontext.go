package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"im-chat/easy-chat/apps/im/rpc/imclient"
	"im-chat/easy-chat/apps/social/api/internal/config"
	"im-chat/easy-chat/apps/social/rpc/socialclient"
	"im-chat/easy-chat/apps/user/rpc/userclient"
	"im-chat/easy-chat/pkg/interceptor"
	"im-chat/easy-chat/pkg/middleware"
)

type ServiceContext struct {
	Config config.Config

	socialclient.Social

	IdempotenceMiddleware rest.Middleware
	LimitMiddleware       rest.Middleware

	userclient.User

	imclient.Im

	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handler,
		LimitMiddleware:       middleware.NewLimitMiddleware(c.Redisx).TokenLimitHandler(1, 100),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc,
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient))),
		User:  userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:    imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		Redis: redis.MustNewRedis(c.Redisx),
	}
}
