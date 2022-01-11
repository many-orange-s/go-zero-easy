package participate

import (
	"context"
	"errors"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"go-zero-easy/service/user/cmd/api/internal/svc"
	"go-zero-easy/service/user/cmd/api/internal/types"
)

type ShowMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) ShowMsgLogic {
	return ShowMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowMsgLogic) ShowMsg() (resp *types.UserMsg, err error) {
	uid := l.ctx.Value("uid")
	usermsg, err := l.svcCtx.UserModel.FindOne(uid.(int64))
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, &errorx.CodeError{Code: errorx.InvalidParam, Msg: errconcrete.UserNotHasMsg}
		} else {
			return nil, &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
		}
	}

	o := &types.UserMsg{
		Account:  usermsg.Account,
		Password: usermsg.Password,
		Name:     usermsg.Name,
		Gender:   usermsg.Gender,
		Address:  usermsg.Address,
		Email:    usermsg.Email,
		Phone:    usermsg.Phone,
	}
	return o, nil
}
