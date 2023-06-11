package logic

import (
	"context"
	"database/sql"

	"shorturl/internal/svc"
	"shorturl/internal/types"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	//  输入betxin.one/2dsa -> 重定向到真实的连接
	// 1. 布隆过滤器
	// 不存在则返回404， 不需要后续处理
	// a 基于内存版本
	// b 基于redis版本
	exist, err := l.svcCtx.Filter.Exists([]byte(req.ShortUrl))
	if err != nil {
		logx.Errorw("svcCtx.Filter.Exists() failed", logx.LogField{Key: "Err", Value: err.Error()})

	}
	if !exist {
		return nil, errorx.NewDefaultError("404")
	}

	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{Valid: true, String: req.ShortUrl})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorx.NewDefaultError("404")
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "Err", Value: err.Error()})
		return nil, err
	}

	// 返回重定向
	return &types.ShowResponse{LongUrl: u.Lurl.String}, nil
}
