package logic

import (
	"context"
	"github.com/tal-tech/go-zero/core/logx"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"go-zero-easy/service/user/cmd/api/internal/svc"
	"go-zero-easy/service/user/cmd/api/internal/types"
	"go-zero-easy/service/user/model"
	"log"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterLogic {
	return RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req types.UserMsg) error {
	account := req.Account
	_, err := l.svcCtx.UserModel.FindOneByAccount(account)
	if err == nil {
		return &errorx.CodeError{
			Code: errorx.InvalidParam,
			Msg:  errconcrete.UserHasExit,
		}
	}

	o := &model.Usermsg{
		Name:     req.Name,
		Gender:   req.Gender,
		Phone:    req.Phone,
		Address:  req.Address,
		Email:    req.Email,
		Account:  req.Account,
		Password: req.Password,
	}
	_, err = l.svcCtx.UserModel.Insert(o)
	if err != nil {
		log.Println("Register insert err:", err)
		return &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errconcrete.InterErr,
		}
	}
	return nil
}
