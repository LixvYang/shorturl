package svc

import (
	"shorturl/internal/config"
	"shorturl/model"
	"shorturl/sequence"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	ShortUrlModel     model.ShortUrlMapModel // short_url_map
	Sequence          sequence.Sequence      // sequence
	ShortUrlBlackList map[string]struct{}
	//bloom filter
	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	// 加载
	m := make(map[string]struct{})
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}

	store := redis.New(c.CacheRedis[0].Host, func(r *redis.Redis) {
		r.Type = redis.NodeType
	})
	filter := bloom.New(store, "bloom_filter", 20*(1<<20))

	return &ServiceContext{
		Config:            c,
		ShortUrlModel:     model.NewShortUrlMapModel(conn, c.CacheRedis),
		Sequence:          sequence.NewMySQL(c.Sequence.DSN),
		ShortUrlBlackList: m,
		Filter:            filter,
	}
}

// 加载已有的短链接数据
func loadDataToBloomFilter() {

}
