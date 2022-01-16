package Members

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"go-zero-easy/service/search/cmd/api/internal/logic/Members"
	"go-zero-easy/service/search/cmd/api/internal/svc"
	"go-zero-easy/service/search/cmd/api/internal/types"
)

func GetBookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BookID
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := Members.NewGetBookLogic(r.Context(), svcCtx)
		err := l.GetBook(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
