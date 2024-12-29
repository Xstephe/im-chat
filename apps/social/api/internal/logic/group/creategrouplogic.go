package group

import (
	"context"
	"im-chat/easy-chat/apps/im/rpc/imclient"
	"im-chat/easy-chat/apps/social/rpc/socialclient"
	"im-chat/easy-chat/pkg/ctxdata"

	"im-chat/easy-chat/apps/social/api/internal/svc"
	"im-chat/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateGroup 创建一个新的群组并建立会话
//
// 功能描述:
//   - 从上下文中获取当前用户ID（群组创建者）
//   - 调用服务层接口创建群组，并获取群组ID
//   - 如果群组创建成功，则建立群组会话
//
// 参数:
//   - req: `*types.GroupCreateReq` 类型，包含群组创建所需的信息（群组名称、图标等）
//
// 返回值:
//   - `*types.GroupCreateResp`: 响应对象，当前未使用
//   - `error`: 如果在创建群组或建立会话过程中发生错误，则返回相应的错误信息
func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUId(l.ctx)

	// 创建群
	res, err := l.svcCtx.Social.GroupCreate(l.ctx, &socialclient.GroupCreateReq{
		Name:       req.Name,
		Icon:       req.Icon,
		CreatorUid: uid,
	})
	if err != nil {
		return nil, err
	}

	if res.Id == "" {
		return nil, nil
	}

	//建立会话
	_, err = l.svcCtx.Im.CreateGroupConversation(l.ctx, &imclient.CreateGroupConversationReq{
		GroupId:  res.Id,
		CreateId: uid,
	})
	return nil, err
}
