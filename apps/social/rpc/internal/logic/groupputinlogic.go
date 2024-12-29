package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"im-chat/easy-chat/apps/social/socialmodels"
	"im-chat/easy-chat/pkg/constants"
	"im-chat/easy-chat/pkg/xerr"
	"time"

	"im-chat/easy-chat/apps/social/rpc/internal/svc"
	"im-chat/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutinLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinLogic {
	return &GroupPutinLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutin 处理群组加入请求
//
// 功能描述:
//   - 根据请求的类型（用户申请、群成员邀请、管理员/创建者邀请）处理群组加入请求。
//   - 如果群组不需要验证，直接将用户加入群组；否则，根据请求的来源和群组角色处理请求。
//
// 参数:
//   - in: `*social.GroupPutinReq` 类型，包含用户ID、群组ID、请求消息、请求时间、加入来源等信息。
//
// 返回值:
//   - `*social.GroupPutinResp`: 处理结果的响应，可能包含群组ID。
//   - `error`: 处理过程中发生的错误。
func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {
	// todo: add your logic here and delete this line
	// todo: add your logic here and delete this line
	//  1. 普通用户申请 ： 如果群无验证直接进入
	//  2. 群成员邀请： 如果群无验证直接进入
	//  3. 群管理员/群创建者邀请：直接进入群

	var (
		inviteGroupMember *socialmodels.GroupMembers
		userGroupMember   *socialmodels.GroupMembers
		groupInfo         *socialmodels.Groups

		err error
	)

	userGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.ReqId, in.GroupId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and  req id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if userGroupMember != nil {
		return &social.GroupPutinResp{}, nil
	}

	groupReq, err := l.svcCtx.GroupRequestsModel.FindByGroupIdAndReqId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group req by groud id and user id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if groupReq != nil {
		return &social.GroupPutinResp{}, nil
	}

	groupReq = &socialmodels.GroupRequests{
		ReqId:   in.ReqId,
		GroupId: in.GroupId,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: sql.NullTime{
			Time:  time.Unix(in.ReqTime, 0),
			Valid: true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	}

	createGroupMember := func() {
		if err != nil {
			return
		}
		err = l.createGroupMember(in)
	}

	groupInfo, err = l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group by groud id err %v, req %v", err, in.GroupId)
	}

	// 验证是否要验证
	if !groupInfo.IsVerify {
		// 不需要
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}

		return l.createGroupReq(groupReq, true)
	}

	// 验证进群方式
	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		// 申请
		return l.createGroupReq(groupReq, false)
	}

	inviteGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.InviterUid, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and user id err %v, req %v",
			in.InviterUid, in.GroupId)
	}

	if constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.CreatorGroupRoleLevel || constants.
		GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.ManagerGroupRoleLevel {
		// 是管理者或创建者邀请
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		groupReq.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		return l.createGroupReq(groupReq, true)
	}
	return l.createGroupReq(groupReq, false)
}

func (l *GroupPutinLogic) createGroupReq(groupReq *socialmodels.GroupRequests, isPass bool) (*social.GroupPutinResp, error) {

	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, groupReq)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert group req err %v req %v", err, groupReq)
	}

	if isPass {
		return &social.GroupPutinResp{GroupId: groupReq.GroupId}, nil
	}

	return &social.GroupPutinResp{}, nil
}

func (l *GroupPutinLogic) createGroupMember(in *social.GroupPutinReq) error {
	groupMember := &socialmodels.GroupMembers{
		GroupId:   in.GroupId,
		UserId:    in.ReqId,
		RoleLevel: int64(constants.AtLargeGroupRoleLevel),
		OperatorUid: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
	}
	_, err := l.svcCtx.GroupMembersModel.Insert(l.ctx, nil, groupMember)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
	}

	return nil
}
