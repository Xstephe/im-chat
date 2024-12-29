package group

import (
	"context"
	"github.com/jinzhu/copier"
	"im-chat/easy-chat/apps/social/rpc/socialclient"

	"im-chat/easy-chat/apps/social/api/internal/svc"
	"im-chat/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutInList 查询未处理的群组加入请求列表
//
// 功能描述:
//   - 根据提供的群组ID从服务层获取所有未处理的群组加入请求
//   - 将获取到的请求转换为响应格式并返回
//
// 参数:
//   - req: `*types.GroupPutInListRep` 类型，包含群组ID，用于查询相关的群组请求
//
// 返回值:
//   - `*types.GroupPutInListResp`: 包含未处理的群组请求列表
//   - `error`: 如果在处理过程中发生错误，则返回相应的错误信息
func (l *GroupPutInListLogic) GroupPutInList(req *types.GroupPutInListRep) (resp *types.GroupPutInListResp, err error) {
	// todo: add your logic here and delete this line
	list, err := l.svcCtx.Social.GroupPutinList(l.ctx, &socialclient.GroupPutinListReq{
		GroupId: req.GroupId,
	})

	var respList []*types.GroupRequests
	copier.Copy(&respList, list.List)

	return &types.GroupPutInListResp{List: respList}, nil
	return
}
