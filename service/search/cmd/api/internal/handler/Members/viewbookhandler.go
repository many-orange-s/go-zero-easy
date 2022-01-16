package Members

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"go-zero-easy/service/search/cmd/api/internal/logic/Members"
	"go-zero-easy/service/search/cmd/api/internal/svc"
)

func ViewBookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := Members.NewViewBookLogic(r.Context(), svcCtx)
		resp, err := l.ViewBook()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
