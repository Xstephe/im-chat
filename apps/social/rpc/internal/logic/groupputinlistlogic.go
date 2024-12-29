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

type GroupPutinListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinListLogic {
	return &GroupPutinListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutinList 查询未处理的群组加入请求列表
//
// 功能描述:
//   - 根据群组ID获取所有未处理的群组加入请求
//   - 将这些请求转换为响应格式并返回
//
// 参数:
//   - in: `*social.GroupPutinListReq` 类型，包含群组ID，用于查询相关的群组请求
//
// 返回值:
//   - `*social.GroupPutinListResp`: 包含未处理的群组请求列表
//   - `error`: 如果在处理过程中发生错误，则返回相应的错误信息
func (l *GroupPutinListLogic) GroupPutinList(in *social.GroupPutinListReq) (*social.GroupPutinListResp, error) {
	// todo: add your logic here and delete this line
	groupReqs, err := l.svcCtx.GroupRequestsModel.ListNoHandler(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group req err %v req %v", err, in.GroupId)
	}

	var respList []*social.GroupRequests
	copier.Copy(&respList, groupReqs)

	return &social.GroupPutinListResp{
		List: respList,
	}, nil
	return &social.GroupPutinListResp{}, nil
}
