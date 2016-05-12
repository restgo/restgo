package restgo

import (
	"net/url"
	"regexp"
	"strings"
)

type layer struct {
	pathRegexp *regexp.Regexp
	path       string
	method     string
	handle     HTTPHandler
	isStatic   bool
	isEnd      bool

	route *Route
}

func newLayer(path string, handle HTTPHandler, end bool) *layer {
	l := &layer{
		handle: handle,
	}
	l.isEnd = end
	l.path = path
	l.pathRegexp, l.isStatic = path2Regexp(path, end)

	return l
}

func (this *layer) handleRequest(ctx *Context, next Next) {
	if this.handle != nil {
		this.handle(ctx, next)
	} else {
		next(nil) // ignore empty handle
	}
}

func (this *layer) match(path string) (url.Values, bool) {
	urlParams := make(url.Values)
	if this.isStatic {
		var match bool
		if this.isEnd == true {
			match = (this.path == path)
		} else {
			match = strings.HasPrefix(path, this.path)
		}

		return urlParams, match
	}

	matches := this.pathRegexp.FindAllStringSubmatch(path, -1)
	if matches != nil {
		names := this.pathRegexp.SubexpNames()
		for i := 1; i < len(names); i++ {
			name := names[i]
			value := matches[0][i]
			if len(name) > 0 && len(value) > 0 {
				urlParams.Set(name, value)
			}
		}
		return urlParams, true
	}

	return urlParams, false
}

// append url params to query string, you can get it by calling ctx.URI().QueryString()
// params := ctx.URI().QueryString()
func (this *layer) registerParamsAsQuery(ctx *Context, urlParams url.Values) {
	query := string(ctx.URI().QueryString())

	for k, v := range urlParams {
		query += ("&" + k + "=" + strings.Join(v, ";"))
	}
	ctx.URI().SetQueryString(query)
}
