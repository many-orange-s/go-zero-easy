package logic

import (
	"context"
	"errors"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"go-zero-easy/service/user/cmd/rpc/internal/svc"
	"go-zero-easy/service/user/cmd/rpc/user"
	"google.golang.org/grpc/codes"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUserMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMsgLogic {
	return &GetUserMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserMsgLogic) GetUserMsg(in *user.UserId) (*user.UserMsg, error) {
	uid := in.Uid

	if uid < 0 {
		return nil, errorx.HandlerRpcError(codes.InvalidArgument, "GetUserMsg uid", errconcrete.RpcUidInvalid)
	}

	userMsg, err := l.svcCtx.UserModel.FindOne(uid)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, errorx.HandlerRpcError(codes.NotFound, "GetUserMsg FindOne", errconcrete.RpcUidNotFound)
		} else {
			return nil, errorx.HandlerRpcError(codes.Internal, "GetUserMsg FindOne", errconcrete.RpcInterErr)
		}
	}

	o := &user.UserMsg{
		Name:    userMsg.Name,
		Gender:  userMsg.Gender,
		Phone:   userMsg.Phone,
		Address: userMsg.Address,
		Email:   userMsg.Email,
	}
	return o, nil
}
