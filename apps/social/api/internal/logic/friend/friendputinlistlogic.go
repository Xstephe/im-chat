package friend

import (
	"context"
	"github.com/jinzhu/copier"
	"im-chat/easy-chat/apps/social/rpc/socialclient"
	"im-chat/easy-chat/pkg/ctxdata"

	"im-chat/easy-chat/apps/social/api/internal/svc"
	"im-chat/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendPutInList 获取用户的好友申请列表
//
// 功能描述:
//   - 从上下文中获取当前用户ID
//   - 调用服务层接口获取当前用户的好友申请列表
//   - 将获取的好友申请列表转换为响应对象并返回
//
// 参数:
//   - req: `*types.FriendPutInListReq` 类型，包含获取好友申请列表的请求信息
//   - `UserId`: 当前用户ID，表示请求获取好友申请列表的用户
//
// 返回值:
//   - `*types.FriendPutInListResp`: 响应对象，包含当前用户的好友申请列表
//   - `List`: 好友申请列表，包含所有待处理的好友申请记录
//   - `error`: 如果在获取好友申请列表过程中发生错误，则返回相应的错误信息
func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	// todo: add your logic here and delete this line
	list, err := l.svcCtx.Social.FriendPutInList(l.ctx, &socialclient.FriendPutInListReq{
		UserId: ctxdata.GetUId(l.ctx),
	})
	if err != nil {
		return nil, err
	}

	var respList []*types.FriendRequests
	copier.Copy(&respList, list.List)
	return &types.FriendPutInListResp{List: respList}, nil
}
