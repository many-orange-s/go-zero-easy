package participate

import (
	"context"
	"encoding/json"
	"fmt"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"go-zero-easy/service/user/cmd/api/internal/svc"
	"go-zero-easy/service/user/cmd/api/internal/types"
	"go-zero-easy/service/user/model"
	"log"

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
	uidNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	uid, err := uidNumber.Int64()
	if err != nil {
		log.Println(err)
		return &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
	}
	log.Println(uid)

	if uid < 0 {
		return &errorx.CodeError{Code: errorx.InvalidParam, Msg: errconcrete.UserUidValid}
	}

	_, err = l.svcCtx.UserModel.FindOne(uid)
	if err != nil {
		return &errorx.CodeError{Code: errorx.InvalidParam, Msg: errconcrete.UserNotHasMsg}
	}

	o := &model.Usermsg{
		Uid:      uid,
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
		log.Println("modifyInformation update err", err)
		return &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
	}
	return nil
}
