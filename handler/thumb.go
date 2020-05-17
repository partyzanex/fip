package handler

import (
	"github.com/nfnt/resize"
	"github.com/valyala/fasthttp"
)

func (h *Handler) Resize(ctx *fasthttp.RequestCtx) {
	h.resize(ctx, resize.Resize)
}

func (h *Handler) Thumb(ctx *fasthttp.RequestCtx) {
	h.resize(ctx, resize.Thumbnail)
}
