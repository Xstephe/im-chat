// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"im-chat/easy-chat/apps/user/rpc/internal/logic"
	"im-chat/easy-chat/apps/user/rpc/internal/svc"
	"im-chat/easy-chat/apps/user/rpc/user"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) Ping(ctx context.Context, in *user.Request) (*user.Response, error) {
	l := logic.NewPingLogic(ctx, s.svcCtx)
	return l.Ping(in)
}

func (s *UserServer) Login(ctx context.Context, in *user.LoginReq) (*user.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

func (s *UserServer) Register(ctx context.Context, in *user.RegisterReq) (*user.RegisterResp, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

func (s *UserServer) GetUserInfo(ctx context.Context, in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	l := logic.NewGetUserInfoLogic(ctx, s.svcCtx)
	return l.GetUserInfo(in)
}

func (s *UserServer) FindUser(ctx context.Context, in *user.FindUserReq) (*user.FindUserResp, error) {
	l := logic.NewFindUserLogic(ctx, s.svcCtx)
	return l.FindUser(in)
}
