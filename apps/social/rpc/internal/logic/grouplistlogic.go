package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"im-chat/easy-chat/pkg/xerr"

	"im-chat/easy-chat/apps/social/rpc/internal/svc"
	"im-chat/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupList 获取用户所在的群组列表
//
// 功能描述:
//   - 根据用户ID查询用户所在的群组成员列表，然后获取这些群组的详细信息，并返回群组列表。
//
// 参数:
//   - in: `social.GroupListReq` 类型，包含用户ID，用于查询用户所在的群组。
//
// 返回值:
//   - `*social.GroupListResp`: 包含用户所在的群组列表的响应对象。
//   - `error`: 如果在获取群组列表过程中发生错误，则返回相应的错误信息。
func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	// todo: add your logic here and delete this line

	userGroup, err := l.svcCtx.GroupMembersModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err %v req %v", err, in.UserId)
	}
	if len(userGroup) == 0 {
		return &social.GroupListResp{}, nil
	}

	ids := make([]string, 0, len(userGroup))
	for _, v := range userGroup {
		ids = append(ids, v.GroupId)
	}

	groups, err := l.svcCtx.GroupsModel.ListByGroupIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group err %v req %v", err, ids)
	}

	var respList []*social.Groups
	copier.Copy(&respList, &groups)

	return &social.GroupListResp{
		List: respList,
	}, nil

	return &social.GroupListResp{}, nil
}
