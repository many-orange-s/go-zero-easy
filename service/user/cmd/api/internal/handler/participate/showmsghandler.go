package participate

import (
	"go-zero-easy/service/user/cmd/api/internal/logic/participate"
	"go-zero-easy/service/user/cmd/api/internal/svc"
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func ShowMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := participate.NewShowMsgLogic(r.Context(), svcCtx)
		resp, err := l.ShowMsg()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
