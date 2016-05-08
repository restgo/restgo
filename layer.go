package grest

import (
	"regexp"
	"net/http"
	"net/url"
	"strings"
)

type layer struct {
	pathRegexp *regexp.Regexp
	path       string
	method     string
	handle     HTTPHandler
	isStatic   bool
	isEnd      bool

	route      *Route
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

func (this *layer) handleRequest(res http.ResponseWriter, req *http.Request, next Next) {
	this.handle.HTTPHandle(res, req, next)
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

// append url params to query string, you can get it by calling req.URL.Query()
// params := req.URL.Query()
// value := params.Get(name)
func (this *layer) registerParamsAsQuery(req *http.Request, urlParams url.Values) {
	query := req.URL.Query()
	for k, v := range urlParams {
		query[k] = v
	}

	req.URL.RawQuery = query.Encode()
}

