package urltool

import (
	"net/url"
	"path"

	"github.com/zeromicro/go-zero/core/logx"
)

func GetBasePath(targetUrl string) (string, error) {
	myUrl, err := url.Parse(targetUrl)
	if err != nil {
		logx.Errorw("url.Parse failed", logx.LogField{Key: "lurl", Value: targetUrl}, logx.LogField{Key: "err", Value: err.Error()})
		return "", err
	}

	return path.Base(myUrl.Path), nil
}
