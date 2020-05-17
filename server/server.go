package server

import (
	"github.com/fasthttp/router"
	"github.com/partyzanex/fip/env"
	"github.com/partyzanex/fip/handler"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"time"
)

func Run(config *env.Config) error {
	h := handler.Handler{
		SourceDir: config.Source,
		CacheDir:  config.Cache,
	}

	r := router.New()
	r.GET("/s/{path}", WithLogger(h.GetSource))
	r.GET("/r/{width}/{height}/{path}", WithLogger(h.Resize))
	r.GET("/t/{width}/{height}/{path}", WithLogger(h.Thumb))

	logrus.Printf("started on %s", config.Address)

	return fasthttp.ListenAndServe(config.Address, r.Handler)
}

func WithLogger(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logrus.WithFields(logrus.Fields{
			"URL":  ctx.Request.URI(),
			"Date": time.Now().String(),
		}).Debug("debug request")

		h(ctx)
	}
}
