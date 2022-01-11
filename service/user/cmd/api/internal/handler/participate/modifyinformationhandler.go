package participate

import (
	"go-zero-easy/service/user/cmd/api/internal/logic/participate"
	"go-zero-easy/service/user/cmd/api/internal/svc"
	"go-zero-easy/service/user/cmd/api/internal/types"
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func ModifyInformationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserMsg
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := participate.NewModifyInformationLogic(r.Context(), svcCtx)
		err := l.ModifyInformation(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
