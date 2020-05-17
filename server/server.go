package server

import (
	"github.com/fasthttp/router"
	"github.com/partyzanex/fip/env"
	"github.com/partyzanex/fip/handler"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func Run(config *env.Config) error {
	h := handler.Handler{
		SourceDir: config.Source,
		CacheDir:  config.Cache,
	}

	r := router.New()
	r.GET("/s/{path}", h.GetSource)
	r.GET("/r/{width}/{height}/{path}", h.Resize)
	r.GET("/t/{width}/{height}/{path}", h.Thumb)

	logrus.Printf("started on %s", config.Address)

	return fasthttp.ListenAndServe(config.Address, r.Handler)
}
