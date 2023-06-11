package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shorturl/internal/svc"
	"shorturl/internal/types"
	"shorturl/model"
	"shorturl/pkg/base62"
	"shorturl/pkg/connect"
	"shorturl/pkg/md5"
	"shorturl/pkg/urltool"

	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConvertLogic) Convert(req *types.ConvertRequset) (resp *types.ConvertResponse, err error) {
	// 1. 校验输入的数据
	// 1.1 数据不能为空 使用validator包来校验
	// 1.2 输入的长链接必须能请求通
	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("url get error!")
	}
	// 1.3 判断之前是否已经转链过（数据库中是否已存在该长链接）
	// 1.3.1 给长链接生成md5
	// 1.3.2 拿md5去数据库中查询是否存在
	md5Value := md5.Sum([]byte(req.LongUrl))
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, errorx.NewDefaultError("该链接已经被转链: " + u.Surl.String)
		}
		logx.Errorw("ShortUrlModel.FindOneByMd5 failed", logx.LogField{Key: "Err", Value: err.Error()})
		return nil, err
	}
	// 1.4 输入的不能是一个短链接
	// 输入的是一个完整的 url，
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("urltool.GetBasePath failed", logx.LogField{Key: "lurl", Value: req.LongUrl})
		return nil, err
	}
	u, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, errorx.NewDefaultError("该链接已经被转链: " + u.Surl.String)
		}
		logx.Errorw("ShortUrlModel.FindOneByMd5 failed", logx.LogField{Key: "Err", Value: err.Error()})
		return nil, err
	}

	var short string
	for {
		// 2. 取号
		// 每来一个REPLACE INTO 语句往sequence 表插入一条数据，并且取出主键id作为号码
		seq, err := l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw("Sequence.Next() failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, err
		}
		fmt.Println(seq)

		// 3. 号码转短链
		// 3.1 安全性
		// 3.2 短域名避免某些特殊词
		short = base62.Int2String(seq)
		fmt.Printf("short:%v\n", short)
		// 4. 存储长链接短链接映射关系
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			break // 生成不在黑名单里的短链接就跳出for循环
		}
	}
	//4.2 将生成的短链接加入布隆过滤器中
	if err = l.svcCtx.Filter.Add([]byte(short)); err != nil {
		logx.Errorw("Filter.Add failed", logx.LogField{Key: "err", Value: err.Error()})
	}

	// 5. 返回响应
	if _, err := l.svcCtx.ShortUrlModel.Insert(l.ctx, &model.ShortUrlMap{
		Lurl: sql.NullString{String: req.LongUrl, Valid: true},
		Md5:  sql.NullString{String: md5Value, Valid: true},
		Surl: sql.NullString{String: short, Valid: true},
	}); err != nil {
		logx.Errorw("Sequence.Next() failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	shortUrl := l.svcCtx.Config.ShortDomain + "/" + short

	return &types.ConvertResponse{ShortUrl: shortUrl}, nil
}
