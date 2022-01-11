package participate

import (
	"context"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"go-zero-easy/service/user/cmd/api/internal/svc"
	"go-zero-easy/service/user/cmd/api/internal/types"
	"go-zero-easy/service/user/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyInformationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyInformationLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyInformationLogic {
	return ModifyInformationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyInformationLogic) ModifyInformation(req types.UserMsg) error {
	uid := l.ctx.Value("uid")
	_, err := l.svcCtx.UserModel.FindOne(uid.(int64))
	if err != nil {
		return &errorx.CodeError{Code: errorx.InvalidParam, Msg: errconcrete.UserNotHasMsg}
	}

	o := &model.Usermsg{
		Uid:      uid.(int64),
		Account:  req.Account,
		Password: req.Password,
		Name:     req.Name,
		Gender:   req.Gender,
		Phone:    req.Phone,
		Address:  req.Address,
		Email:    req.Email,
	}
	err = l.svcCtx.UserModel.Update(o)
	if err != nil {
		return &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
	}
	return nil
}
