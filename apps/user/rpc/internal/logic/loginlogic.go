package logic

import (
	"context"
	"github.com/pkg/errors"
	"im-chat/easy-chat/apps/user/models"
	"im-chat/easy-chat/pkg/ctxdata"
	"im-chat/easy-chat/pkg/encrypt"
	"im-chat/easy-chat/pkg/xerr"
	"time"

	"im-chat/easy-chat/apps/user/rpc/internal/svc"
	"im-chat/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号码没有注册")
	ErrUserPwdError     = xerr.New(xerr.SERVER_COMMON_ERROR, "密码是错误的")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 处理用户登录请求。
//
// 功能描述:
//   - 验证用户是否已经注册。
//   - 验证用户提供的密码是否正确。
//   - 如果验证通过，为用户生成 JWT 令牌并返回。
//
// 参数:
//   - in: 包含用户登录信息的请求结构体。
//
// 返回值:
//   - *user.LoginResp: 包含用户ID、JWT令牌、过期时间和用户信息的响应结构体。
//   - error: 如果登录过程中出现错误，则返回相应的错误信息。
func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {

	//1.验证用户是否注册过，根据手机号码去查询
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v, req %v", in.Phone, in.Phone)
	}

	//密码验证
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	//生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", in.Phone)
	}

	return &user.LoginResp{
		Id:     userEntity.Id,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
