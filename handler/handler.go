package handler

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/amalfra/etag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	pathReplacer = strings.NewReplacer("./", "", "../", "")

	jpegHeader   = []byte("image/jpeg")
	pngHeader    = []byte("image/png")
	cacheControl = []byte("public, max-age=3600")

	jpegOptions = &jpeg.Options{Quality: 80}

	headerETag         = []byte(fasthttp.HeaderETag)
	headerLastModified = []byte(fasthttp.HeaderLastModified)
	headerCacheControl = []byte(fasthttp.HeaderCacheControl)
	headerExpires      = []byte(fasthttp.HeaderExpires)
)

type Handler struct {
	SourceDir string
	CacheDir  string
}

func (Handler) internalErr(ctx *fasthttp.RequestCtx, err error) {
	logrus.Error(err)
	ctx.SetContentType("text/plain")
	ctx.Error(err.Error(), http.StatusInternalServerError)
}

func (Handler) badRequestErr(ctx *fasthttp.RequestCtx, err error) {
	logrus.Error(err)
	ctx.SetContentType("text/plain")
	ctx.Error(err.Error(), http.StatusBadRequest)
}

func (Handler) close(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		logrus.Error(err)
	}
}

func (Handler) encodeImage(w io.Writer, img image.Image, ext string) error {
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(w, img, jpegOptions)
	case ".png":
		return png.Encode(w, img)
	case ".gif":
		return gif.Encode(w, img, nil)
	default:
		return errors.Errorf("unknown image format %s", ext)
	}
}

func (Handler) setContentType(ctx *fasthttp.RequestCtx, ext string) {
	switch ext {
	case ".jpg", ".jpeg":
		ctx.Response.Header.SetContentTypeBytes(jpegHeader)
	case ".png":
		ctx.Response.Header.SetContentTypeBytes(pngHeader)
	}
}

func (Handler) setHeaders(ctx *fasthttp.RequestCtx, info os.FileInfo) {
	ctx.Response.Header.SetBytesK(headerETag, etag.Generate(info.Name(), true))
	ctx.Response.Header.SetBytesK(headerLastModified, info.ModTime().Format(time.RFC1123))
	ctx.Response.Header.SetBytesKV(headerCacheControl, cacheControl)
	ctx.Response.Header.SetBytesK(headerExpires, time.Now().Add(3600).Format(time.RFC1123))
}

func (Handler) getSource(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "opening source failed")
	}

	return file, nil
}

func (Handler) sourceIsExists(path string) (os.FileInfo, bool, error) {
	info, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, false, errors.Wrap(err, "getting file info failed")
	}

	return info, info != nil, nil
}

func (Handler) sanitizePath(path string) string {
	return pathReplacer.Replace(path)
}
