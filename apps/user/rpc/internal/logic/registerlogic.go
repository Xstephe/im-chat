package logic

import (
	"context"
	"database/sql"
	"errors"
	"im-chat/easy-chat/apps/user/models"
	"im-chat/easy-chat/pkg/ctxdata"
	"im-chat/easy-chat/pkg/encrypt"
	"im-chat/easy-chat/pkg/wuid"
	"time"

	"im-chat/easy-chat/apps/user/rpc/internal/svc"
	"im-chat/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegister = errors.New("手机号码已经注册")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Register 处理用户注册请求。
//
// 功能描述:
//   - 验证用户是否已经注册。
//   - 如果用户未注册，创建一个新的用户并将其信息存入数据库。
//   - 为用户生成 JWT 令牌并返回。
//
// 参数:
//   - in: 包含用户注册信息的请求结构体。
//
// 返回值:
//   - *user.RegisterResp: 包含生成的 JWT 令牌和过期时间的响应结构体。
//   - error: 如果注册过程中出现错误，则返回相应的错误信息。
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {

	//1. 验证用户是否注册，根据手机号进行验证
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}

	if userEntity != nil {
		return nil, ErrPhoneIsRegister
	}

	//定义用户数据
	userEntity = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	if len(in.Password) > 0 {
		//给密码加密
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		//设置密码值
		userEntity.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}

	//插入数据
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, err
	}

	//生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
