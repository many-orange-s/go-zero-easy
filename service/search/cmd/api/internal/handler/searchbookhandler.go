package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"go-zero-easy/service/search/cmd/api/internal/logic"
	"go-zero-easy/service/search/cmd/api/internal/svc"
	"go-zero-easy/service/search/cmd/api/internal/types"
)

func SearchBookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BookName
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSearchBookLogic(r.Context(), svcCtx)
		resp, err := l.SearchBook(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
