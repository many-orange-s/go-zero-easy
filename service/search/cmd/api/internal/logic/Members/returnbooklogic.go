package Members

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	redis_lock "go-zero-easy/service/search/cmd/api/internal/redis-lock"
	"log"

	"go-zero-easy/service/search/cmd/api/internal/svc"
	"go-zero-easy/service/search/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ReturnBookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReturnBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) ReturnBookLogic {
	return ReturnBookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReturnBookLogic) ReturnBook(req types.BookID) error {
	uidNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	uid, err := uidNumber.Int64()
	if err != nil {
		log.Println("ReturnBook getuid err", err)
		return &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errorx.InternalInfo,
		}
	}
	bookid := req.Id

	rendmsg, err := l.svcCtx.LendModel.FindOneByBookidUid(bookid, uid)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return &errorx.CodeError{
				Code: errorx.InvalidParam,
				Msg:  errconcrete.BookNotRendByYou,
			}
		} else {
			log.Println("ReturnBook FindOneByBookidUid err", err)
			return &errorx.CodeError{
				Code: errorx.SystemBusy,
				Msg:  errorx.InternalInfo,
			}
		}
	}

	err = l.svcCtx.LendModel.Delete(rendmsg.Id)
	if err != nil {
		log.Println("ReturnBook .Delete err", err)
		return &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errorx.InternalInfo,
		}
	}

	mutex := redis_lock.ReturnLock(rendmsg.Bookid)
	if err = mutex.Lock(); err != nil {
		log.Println("GetBook getLock err", err)
		return &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errconcrete.InterErr,
		}
	}
	defer mutex.Unlock()

	bookmsg, err := l.svcCtx.BookModel.FindOneByBookId(rendmsg.Bookid)
	if err != nil {
		log.Println("GetBook FindOneByBookId err", err)
		return &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errconcrete.InterErr,
		}
	}

	bookmsg.Count += 1
	err = l.svcCtx.BookModel.Update(bookmsg)
	if err != nil {
		log.Println("GetBook Update err", err)
		return &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errconcrete.InterErr,
		}
	}

	return nil
}
