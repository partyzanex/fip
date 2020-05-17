package handler

import (
	"image"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Resize func(maxWidth, maxHeight uint, img image.Image, interp resize.InterpolationFunction) image.Image

func (h *Handler) resize(ctx *fasthttp.RequestCtx, r Resize) {
	width, err := strconv.Atoi(ctx.UserValue("width").(string))
	if err != nil {
		h.badRequestErr(ctx, err)
		return
	}

	height, err := strconv.Atoi(ctx.UserValue("height").(string))
	if err != nil {
		h.badRequestErr(ctx, err)
		return
	}

	if width == 0 && height == 0 {
		h.badRequestErr(ctx, errors.New("empty size"))
		return
	}

	if width == 0 {
		width = height
	}

	if height == 0 {
		height = width
	}

	path := ctx.UserValue("path").(string)
	path = filepath.Join(h.SourceDir, h.sanitizePath(path))

	cacheFile := filepath.Join(h.CacheDir, string(ctx.Request.URI().Path()))
	info, exists, err := h.sourceIsExists(cacheFile)
	if err != nil {
		h.internalErr(ctx, err)
		return
	}

	if exists {
		reader, err := h.getSource(cacheFile)
		if err != nil {
			h.internalErr(ctx, err)
			return
		}

		defer func() {
			err := reader.Close()
			if err != nil {
				logrus.Error(err)
			}
		}()

		h.setContentType(ctx, strings.ToLower(filepath.Ext(cacheFile)))
		h.setHeaders(ctx, info)

		_, err = io.Copy(ctx, reader)
		if err != nil {
			h.internalErr(ctx, err)
			return
		}

		return
	}

	info, exists, err = h.sourceIsExists(path)
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

	img, ext, err := image.Decode(reader)
	if err != nil {
		h.internalErr(ctx, err)
		return
	}

	ext = "." + ext

	cachePath := filepath.Dir(cacheFile)

	err = os.MkdirAll(cachePath, 0777)
	if err != nil {
		h.internalErr(ctx, err)
		return
	}

	file, err := os.Create(cacheFile)
	if err != nil {
		h.internalErr(ctx, err)
		return
	}

	defer h.close(file)

	img = r(uint(width), uint(height), img, resize.MitchellNetravali)

	err = h.encodeImage(file, img, ext)
	if err != nil {
		h.internalErr(ctx, err)
		return
	}

	h.setContentType(ctx, ext)
	h.setHeaders(ctx, info)

	err = h.encodeImage(ctx, img, ext)
	if err != nil {
		h.internalErr(ctx, err)
		return
	}
}
