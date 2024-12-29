package group

import (
	"context"
	"github.com/jinzhu/copier"
	"im-chat/easy-chat/apps/social/rpc/socialclient"
	"im-chat/easy-chat/pkg/ctxdata"

	"im-chat/easy-chat/apps/social/api/internal/svc"
	"im-chat/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupList 获取用户所在的群组列表
//
// 功能描述:
//   - 从上下文中获取当前用户ID
//   - 调用服务层接口获取用户所在的群组列表
//   - 将获取的群组信息转换为响应格式并返回
//
// 参数:
//   - req: `*types.GroupListReq` 类型，包含请求获取群组列表的信息（当前未使用）
//
// 返回值:
//   - `*types.GroupListResp`: 包含用户所在群组列表的响应对象
//   - `error`: 如果获取群组列表过程中发生错误，则返回相应的错误信息
func (l *GroupListLogic) GroupList(req *types.GroupListRep) (resp *types.GroupListResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUId(l.ctx)
	list, err := l.svcCtx.Social.GroupList(l.ctx, &socialclient.GroupListReq{
		UserId: uid,
	})

	if err != nil {
		return nil, err
	}

	var respList []*types.Groups
	copier.Copy(&respList, list.List)

	return &types.GroupListResp{List: respList}, nil
	return
}
