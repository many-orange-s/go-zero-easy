package Members

import (
	"context"
	"encoding/json"
	"fmt"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"log"

	"go-zero-easy/service/search/cmd/api/internal/svc"
	"go-zero-easy/service/search/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ViewHowBookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewViewHowBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) ViewHowBookLogic {
	return ViewHowBookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ViewHowBookLogic) ViewHowBook() (resp *types.Count, err error) {
	uidNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	uid, err := uidNumber.Int64()
	if err != nil {
		log.Println("ViewHowBook getuid err", err)
		return nil, &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errconcrete.InterErr,
		}
	}

	num, err := l.svcCtx.LendModel.SearchHowManyBook(uid)
	if err != nil {
		log.Println("ViewHowBook SearchHowManyBook err : ", err)
		return nil, &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errconcrete.InterErr,
		}
	}
	resp = new(types.Count)
	resp.Num = num
	return
}
