package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/lixvyang/rebetxin-one/common/errorx"

	"shorturl/internal/logic"
	"shorturl/internal/svc"
	"shorturl/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultParamsFailedError())
			return
		}

		// 参数规则校验
		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			logx.Error("validator check failed: ", logx.LogField{Key: "Err", Value: err.Error()})
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShowLogic(r.Context(), svcCtx)
		resp, err := l.Show(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			http.Redirect(w, r, resp.LongUrl, http.StatusFound)
		}
	}
}
