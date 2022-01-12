package participate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"go-zero-easy/service/user/cmd/api/internal/svc"
	"go-zero-easy/service/user/cmd/api/internal/types"
	"log"
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
	// 这里面不能直接从l.ctx.Value("userId")拿出来会报错
	// 要这样子转换
	uidNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	uid, err := uidNumber.Int64()
	if err != nil {
		log.Println(err)
		return nil, &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
	}
	log.Println(uid)

	if uid < 0 {
		return nil, &errorx.CodeError{Code: errorx.InvalidParam, Msg: errconcrete.UserUidValid}
	}

	usermsg, err := l.svcCtx.UserModel.FindOne(uid)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, &errorx.CodeError{Code: errorx.InvalidParam, Msg: errconcrete.UserNotHasMsg}
		} else {
			log.Println("showMsg FindOne err", err)
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
