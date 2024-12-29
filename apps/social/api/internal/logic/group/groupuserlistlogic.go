package group

import (
	"context"
	"im-chat/easy-chat/apps/social/rpc/socialclient"
	"im-chat/easy-chat/apps/user/rpc/userclient"

	"im-chat/easy-chat/apps/social/api/internal/svc"
	"im-chat/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserListLogic {
	return &GroupUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupUserList 处理查询群组成员及其用户信息的请求
//
// 功能描述:
//   - 该方法查询指定群组的所有成员，并获取这些成员的用户信息。
//   - 将群组成员及其对应的用户信息构造为响应返回。
//
// 参数:
//   - req: `*types.GroupUserListReq` 类型，包含查询所需的群组ID。
//
// 返回值:
//   - `*types.GroupUserListResp`: 包含群组成员及其用户信息的响应。
//   - `error`: 处理过程中发生的错误。
func (l *GroupUserListLogic) GroupUserList(req *types.GroupUserListReq) (resp *types.GroupUserListResp, err error) {
	// todo: add your logic here and delete this line
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})

	// 还需要获取用户的信息
	uids := make([]string, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		uids = append(uids, v.UserId)
	}

	// 获取用户信息
	userList, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{Ids: uids})
	if err != nil {
		return nil, err
	}

	userRecords := make(map[string]*userclient.UserEntity, len(userList.User))
	for i, _ := range userList.User {
		userRecords[userList.User[i].Id] = userList.User[i]
	}

	respList := make([]*types.GroupMembers, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {

		member := &types.GroupMembers{
			Id:        int64(v.Id),
			GroupId:   v.GroupId,
			UserId:    v.UserId,
			RoleLevel: int(v.RoleLevel),
		}
		if u, ok := userRecords[v.UserId]; ok {
			member.Nickname = u.Nickname
			member.UserAvatarUrl = u.Avatar
		}
		respList = append(respList, member)
	}

	return &types.GroupUserListResp{List: respList}, err
}
