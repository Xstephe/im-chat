package main

import (
	"flag"
	"fmt"
	"im-chat/easy-chat/apps/social/rpc/internal/config"
	"im-chat/easy-chat/apps/social/rpc/internal/server"
	"im-chat/easy-chat/apps/social/rpc/internal/svc"
	"im-chat/easy-chat/apps/social/rpc/social"
	"im-chat/easy-chat/pkg/interceptor"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		social.RegisterSocialServer(grpcServer, server.NewSocialServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	//s.AddUnaryInterceptors(rpcserver.LogInterceptor, rpcserver.SyncXLimitInterceptor(10))
	s.AddUnaryInterceptors(interceptor.NewIdempotenceServer(interceptor.NewDefaultIdempotent(c.Cache[0].RedisConf)))
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
