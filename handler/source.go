package handler

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

func (h *Handler) GetSource(ctx *fasthttp.RequestCtx) {
	path := ctx.UserValue("path").(string)
	path = filepath.Join(h.SourceDir, h.sanitizePath(path))

	info, exists, err := h.sourceIsExists(path)
	if err != nil {
		h.internalErr(ctx, err)
		return
	}

	if !exists {
		ctx.NotFound()
		return
	}

	reader, err := h.getSource(path)
	if err != nil {
		h.internalErr(ctx, err)
		return
	}

	defer h.close(reader)

	h.setContentType(ctx, strings.ToLower(filepath.Ext(info.Name())))
	h.setHeaders(ctx, info)

	_, err = io.Copy(ctx, reader)
	if err != nil {
		h.internalErr(ctx, err)
		return
	}
}
