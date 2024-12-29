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

type FriendPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutInList 获取用户的好友请求列表
//
// 功能描述:
//   - 从数据库中查询指定用户的所有未处理的好友请求。
//   - 将查询结果转换为统一的响应格式并返回。
//
// 参数:
//   - in: `social.FriendPutInListReq` 类型，包含用户ID，用于查询该用户的好友请求列表。
//
// 返回值:
//   - `*social.FriendPutInListResp`: 包含好友请求列表的响应对象。
//   - `error`: 如果查询过程中发生错误，则返回相应的错误信息。
func (l *FriendPutInListLogic) FriendPutInList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	// todo: add your logic here and delete this line
	friendReqList, err := l.svcCtx.FriendRequestsModel.ListNoHandler(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find list friend req err %v req %v", err, in.UserId)
	}

	var resp []*social.FriendRequests
	copier.Copy(&resp, &friendReqList)

	return &social.FriendPutInListResp{
		List: resp,
	}, nil
	return &social.FriendPutInListResp{}, nil
}
