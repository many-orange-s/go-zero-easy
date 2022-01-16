package Members

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"go-zero-easy/service/search/cmd/api/internal/logic/Members"
	"go-zero-easy/service/search/cmd/api/internal/svc"
)

func ViewHowBookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := Members.NewViewHowBookLogic(r.Context(), svcCtx)
		resp, err := l.ViewHowBook()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
