package logic

import (
	"context"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"log"

	"go-zero-easy/service/search/cmd/api/internal/svc"
	"go-zero-easy/service/search/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SearchBookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) SearchBookLogic {
	return SearchBookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchBookLogic) SearchBook(req types.BookName) (resp *types.BookSet, err error) {
	name := req.Name
	if name == "" {
		log.Println("searchBook get name err", err)
		return nil, &errorx.CodeError{
			Code: errorx.InvalidParam,
			Msg:  errconcrete.BookNameValid,
		}
	}

	bookmsg, err := l.svcCtx.BookModel.SearchByName(name)
	if err != nil {
		log.Println("searchBook SearchByName err", err)
		return nil, &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errconcrete.InterErr,
		}
	}
	if bookmsg == nil {
		return nil, &errorx.CodeError{
			Code: errorx.InvalidParam,
			Msg:  errconcrete.BookNotFount,
		}
	}

	resp = new(types.BookSet)
	resp.BookMsgs = make([]*types.BookMsg, 0, 10)
	for _, msg := range bookmsg {
		o := &types.BookMsg{
			Number: msg.BookId,
			Name:   msg.Name,
			Count:  int(msg.Count),
		}
		resp.BookMsgs = append(resp.BookMsgs, o)
	}
	return
}
