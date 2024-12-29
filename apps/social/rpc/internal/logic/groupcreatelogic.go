package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"im-chat/easy-chat/apps/social/socialmodels"
	"im-chat/easy-chat/pkg/constants"
	"im-chat/easy-chat/pkg/wuid"
	"im-chat/easy-chat/pkg/xerr"

	"im-chat/easy-chat/apps/social/rpc/internal/svc"
	"im-chat/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupCreate 创建一个新的群组
//
// 功能描述:
//   - 根据请求参数创建一个新的群组，并将群组信息和群组成员信息插入到数据库中。
//   - 在插入过程中使用事务来确保数据的一致性和完整性。
//
// 参数:
//   - in: `social.GroupCreateReq` 类型，包含创建群组所需的信息，包括群组名称、图标、创建者ID等。
//
// 返回值:
//   - `*social.GroupCreateResp`: 包含新创建群组的ID的响应对象。
//   - `error`: 如果在创建过程中发生错误，则返回相应的错误信息。
func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// todo: add your logic here and delete this line
	groups := &socialmodels.Groups{
		Id:         wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Name:       in.Name,
		Icon:       in.Icon,
		CreatorUid: in.CreatorUid,
		//IsVerify:   true,
		IsVerify: false,
	}

	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.GroupsModel.Insert(l.ctx, session, groups)

		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group err %v req %v", err, in)
		}

		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, &socialmodels.GroupMembers{
			GroupId:   groups.Id,
			UserId:    in.CreatorUid,
			RoleLevel: int64((constants.CreatorGroupRoleLevel)),
		})
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v req %v", err, in)
		}
		return nil
	})

	return &social.GroupCreateResp{
		Id: groups.Id,
	}, err
}
