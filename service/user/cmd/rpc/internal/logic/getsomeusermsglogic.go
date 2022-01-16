package logic

import (
	"context"
	"errors"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"google.golang.org/grpc/codes"
	"io"

	"go-zero-easy/service/user/cmd/rpc/internal/svc"
	"go-zero-easy/service/user/cmd/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetSomeUserMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSomeUserMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSomeUserMsgLogic {
	return &GetSomeUserMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSomeUserMsgLogic) GetSomeUserMsg(stream user.User_GetSomeUserMsgServer) error {
	for {
		uid, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return errorx.HandlerRpcError(codes.Internal, "GetSomeUserMsg Recv", errconcrete.RpcInterErr)
		}

		userMsg, err := l.svcCtx.UserModel.FindOne(uid.Uid)
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return errorx.HandlerRpcError(codes.NotFound, "GetSomeUserMsg FindOne", errconcrete.RpcUidNotFound)
			} else {
				return errorx.HandlerRpcError(codes.Internal, "GetSomeUserMsg FindOne", errconcrete.RpcInterErr)
			}
		}

		o := &user.UserMsg{
			Name:    userMsg.Name,
			Gender:  userMsg.Gender,
			Phone:   userMsg.Phone,
			Address: userMsg.Address,
			Email:   userMsg.Email,
		}

		err = stream.Send(o)
		if err != nil {
			return errorx.HandlerRpcError(codes.Internal, "GetSomeUserMsg send", errconcrete.RpcInterErr)
		}
	}

	return nil
}
