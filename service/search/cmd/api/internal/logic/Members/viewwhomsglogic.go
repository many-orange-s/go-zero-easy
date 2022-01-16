package Members

import (
	"context"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"go-zero-easy/service/user/cmd/rpc/user"
	"log"

	"go-zero-easy/service/search/cmd/api/internal/svc"
	"go-zero-easy/service/search/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ViewWhoMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewViewWhoMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) ViewWhoMsgLogic {
	return ViewWhoMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ViewWhoMsgLogic) ViewWhoMsg(req types.BookID) (resp *types.PeopleMsg, err error) {
	bookid := req.Id

	nums, err := l.svcCtx.LendModel.SearchWhoLookBook(bookid)
	if err != nil {
		log.Println("ViewWhoMsg SearchWhoLookBook err:", err)
		return nil, &errorx.CodeError{
			Code: errorx.SystemBusy,
			Msg:  errorx.InternalInfo,
		}
	}
	if nums == nil {
		return nil, &errorx.CodeError{
			Code: errorx.InvalidParam,
			Msg:  errconcrete.BookNotRend,
		}
	}

	resp = new(types.PeopleMsg)
	resp.People = make([]*types.Msg, 0, 10)

	if len(nums) == 1 {
		o := user.UserId{
			Uid: *nums[0],
		}
		userMsg, err := l.svcCtx.UserRpc.GetUserMsg(l.ctx, &o)
		if err != nil {
			log.Println("ViewWhoMsg .GetUserMsg err :", err)
			return nil, &errorx.CodeError{
				Code: errorx.SystemBusy,
				Msg:  errorx.InternalInfo,
			}
		}
		u := &types.Msg{
			Name:    userMsg.Name,
			Gender:  userMsg.Gender,
			Email:   userMsg.Email,
			Phone:   userMsg.Phone,
			Address: userMsg.Address,
		}
		resp.People = append(resp.People, u)
	} else {
		stream, err := l.svcCtx.UserRpc.GetSomeUserMsg(l.ctx)
		if err != nil {
			log.Println("ViewWhoMsg GetSomeUserMsg err", err)
			return nil, &errorx.CodeError{
				Code: errorx.SystemBusy,
				Msg:  errorx.InternalInfo,
			}
		}

		for _, num := range nums {
			u := &user.UserId{
				Uid: *num,
			}
			err = stream.Send(u)
			if err != nil {
				log.Println("ViewWhoMsg Send err :", err)
				return nil, &errorx.CodeError{
					Code: errorx.SystemBusy,
					Msg:  errorx.InternalInfo,
				}
			}

			usermsg, err := stream.Recv()
			if err != nil {
				log.Println("ViewWhoMsg Recv err :", err)
				return nil, &errorx.CodeError{
					Code: errorx.SystemBusy,
					Msg:  errorx.InternalInfo,
				}
			}

			m := &types.Msg{
				Name:    usermsg.Name,
				Gender:  usermsg.Gender,
				Email:   usermsg.Email,
				Phone:   usermsg.Phone,
				Address: usermsg.Address,
			}
			resp.People = append(resp.People, m)
		}
		_ = stream.CloseSend()
	}

	return resp, nil
}
