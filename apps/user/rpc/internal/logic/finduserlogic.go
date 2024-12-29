package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"im-chat/easy-chat/apps/user/models"

	"im-chat/easy-chat/apps/user/rpc/internal/svc"
	"im-chat/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FindUser 根据请求参数查找用户。
//
// 功能描述:
//   - 根据不同的请求参数（手机号、用户名、用户ID列表）从数据库中查找用户，可以查找单个或多个用户。
//   - 如果查询成功，将用户实体填充到响应对象中并返回。
//   - 如果查询过程中出现错误，则返回相应的错误信息。
//
// 参数:
//   - in: *user.FindUserReq
//     包含查找用户所需的请求参数（手机号、用户名、用户ID列表）。
//
// 返回值:
//   - *user.FindUserResp: 包含查询到的用户信息的响应对象。
//   - error: 如果查询过程中发生错误，返回相应的错误信息。
func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var (
		userEntitys []*models.Users
		err         error
	)
	if in.Phone != "" {
		userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntitys = append(userEntitys, userEntity)
		}
	} else if in.Name != "" {
		userEntitys, err = l.svcCtx.UsersModel.ListByname(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		userEntitys, err = l.svcCtx.UsersModel.ListByids(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, err
	}

	var resp []*user.UserEntity
	copier.Copy(&resp, userEntitys)
	return &user.FindUserResp{
		User: resp,
	}, nil
}
