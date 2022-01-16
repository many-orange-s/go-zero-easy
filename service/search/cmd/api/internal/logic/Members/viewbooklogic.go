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

type ViewBookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewViewBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) ViewBookLogic {
	return ViewBookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ViewBookLogic) ViewBook() (resp *types.BookSet, err error) {
	uidNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	uid, err := uidNumber.Int64()
	if err != nil {
		log.Println("ViewBook getuid err", err)
		return nil, &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errorx.InternalInfo,
		}
	}

	log.Println(uid)

	bookIds, err := l.svcCtx.LendModel.SearchAllMsgByUid(uid)
	if err != nil {
		log.Println("ViewBook SearchAllMsgByUid err", err)
		return nil, &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errorx.InternalInfo,
		}
	}
	if bookIds == nil {
		return nil, &errorx.CodeError{
			Code: errorx.InvalidParam,
			Msg:  errconcrete.BookNotAnyRend,
		}
	}

	log.Println(bookIds)
	bookmsgs, err := l.svcCtx.BookModel.FindAllMsg(bookIds)
	if err != nil {
		log.Println("ViewBook FindAllMsg err", err)
		return nil, &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errorx.InternalInfo,
		}
	}

	resp = new(types.BookSet)
	resp.BookMsgs = make([]*types.BookMsg, 0, 10)
	for _, bookmsg := range bookmsgs {
		o := &types.BookMsg{
			Number: bookmsg.BookId,
			Name:   bookmsg.Name,
			Count:  int(bookmsg.Count),
		}
		resp.BookMsgs = append(resp.BookMsgs, o)
	}

	return
}
