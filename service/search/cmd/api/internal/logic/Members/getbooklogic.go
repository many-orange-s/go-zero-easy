package Members

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	redis_lock "go-zero-easy/service/search/cmd/api/internal/redis-lock"
	"go-zero-easy/service/search/cmd/api/internal/svc"
	"go-zero-easy/service/search/cmd/api/internal/types"
	"go-zero-easy/service/search/model"
	"log"
)

type GetBookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetBookLogic {
	return GetBookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBookLogic) GetBook(req types.BookID) error {
	bookid := req.Id
	uidNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	uid, err := uidNumber.Int64()

	log.Printf("uid :%v , bookid : %v\n", uid, bookid)

	if err != nil {
		log.Println("GetBook getuid err", err)
		return &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errorx.InternalInfo,
		}
	}

	mutex := redis_lock.ReturnLock(bookid)
	if err = mutex.Lock(); err != nil {
		log.Println("GetBook getLock err", err)
		return &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errconcrete.InterErr,
		}
	}
	defer mutex.Unlock()

	bookmeg, err := l.svcCtx.BookModel.FindOneByBookId(bookid)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return &errorx.CodeError{Code: errorx.InvalidParam, Msg: errconcrete.BookNotFound}
		} else {
			log.Println("GetBook FindOneByBookId err ", err)
			return &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
		}
	}

	if bookmeg.Count <= 0 {
		return &errorx.CodeError{Code: errorx.InvalidParam, Msg: errconcrete.BookNotHasInventory}
	}

	_, err = l.svcCtx.LendModel.FindOneByBookidUid(bookid, uid)
	if err != nil {
		if !errors.Is(err, sqlx.ErrNotFound) {
			log.Println("GetBook FindOneByBookIdUid err : ", err)
			return &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
		}
	} else {
		return &errorx.CodeError{
			Code: errorx.InvalidParam,
			Msg:  errconcrete.BookHasRendByYou,
		}
	}

	o := &model.Bookmsg{
		Id:     bookmeg.Id,
		BookId: bookmeg.BookId,
		Name:   bookmeg.Name,
		Count:  bookmeg.Count - 1,
	}
	err = l.svcCtx.BookModel.Update(o)
	if err != nil {
		log.Println("GetBook Update err :", err)
		return &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
	}

	mutex.Unlock()

	u := &model.Lend{
		Bookid: bookid,
		Uid:    uid,
	}
	_, err = l.svcCtx.LendModel.Insert(u)
	if err != nil {
		log.Println("GetBook Insert err :", err)
		return &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
	}

	return nil
}
