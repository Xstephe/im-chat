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

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutIn 处理添加好友请求
//
// 功能描述:
//   - 检查申请人和目标用户是否已是好友。
//   - 检查是否已有未处理的好友请求。
//   - 如果未找到好友关系或已有未处理的好友请求，则创建新的好友请求记录。
//
// 参数:
//   - in: `social.FriendPutInReq` 类型，包含用户ID、请求用户ID、请求消息和请求时间等信息。
//
// 返回值:
//   - `*social.FriendPutInResp`: 包含处理结果的响应对象。
//   - `error`: 如果发生错误，则返回相应的错误信息。
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: add your logic here and delete this line
	//1.申请人与目标是否是好友关系
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFId(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by uid and fid err %v req %v ", err, in)
	}

	if friends != nil {
		return &social.FriendPutInResp{}, nil
	}

	//2.是否已经有过申请，申请不成功，没有完成
	friendReqs, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by rid and uid err %v req %v ", err, in)
	}
	if friendReqs != nil {
		return &social.FriendPutInResp{}, err
	}

	//3.创建申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest err %v req %v ", err, in)
	}

	return &social.FriendPutInResp{}, nil
}
