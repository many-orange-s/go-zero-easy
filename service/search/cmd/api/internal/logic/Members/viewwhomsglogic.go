package Members

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
