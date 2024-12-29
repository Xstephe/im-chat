package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"im-chat/easy-chat/apps/im/ws/internal/config"
	"im-chat/easy-chat/apps/im/ws/internal/handler"
	"im-chat/easy-chat/apps/im/ws/internal/svc"
	"im-chat/easy-chat/apps/im/ws/websocket"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	//设置日志监听等处理
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		//websocket.WithServerAck(websocket.RigorAck),
		//websocket.WithServerMaxConnectionIdle(10*time.Second),
	)
	defer srv.Stop()

	//加载路由
	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")
	srv.Start()

}
