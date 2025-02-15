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

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendList 获取用户的好友列表
//
// 功能描述:
//   - 根据用户ID从数据库中查询用户的好友列表。
//   - 将查询到的好友列表转换为响应格式并返回。
//
// 参数:
//   - in: `social.FriendListReq` 类型，包含用户ID (`UserId`)。
//
// 返回值:
//   - `*social.FriendListResp`: 包含好友列表的响应对象。
//   - `error`: 如果查询过程中发生错误，则返回相应的错误信息。
func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	// todo: add your logic here and delete this line
	friendsList, err := l.svcCtx.FriendsModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list friend by uid err %v req %v ", err,
			in.UserId)
	}

	var respList []*social.Friends
	copier.Copy(&respList, &friendsList)

	return &social.FriendListResp{
		List: respList,
	}, nil
}
