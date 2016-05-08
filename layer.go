package grester

import (
	"regexp"
	"net/http"
)

type layer struct {
	regexp *regexp.Regexp
	path string
	method string
	handle HTTPHandler
	params map[int] string
	route *Route
}

func newLayer(path string, handle HTTPHandler, end bool) *layer {
	l := &layer{
		handle: handle,
	}
	l.regexp, l.params= Path2Regexp(path, end)

	return l
}

func (this *layer) handleRequest (res http.ResponseWriter, req *http.Request, next Next) {
	this.handle.HTTPHandle(res, req, next)
}


func (this *layer) match (path string) bool{
	if path == "" {
		return false
	}

	m := this.regexp.FindStringSubmatch(path)
	if len(m) == 0 {
		return false
	}

	this.path = m[0]
	return true
}

func (this *layer) registerParamsAsQuery(path string, req *http.Request){

	m := this.regexp.FindStringSubmatch(path)
	if len(m) > 1 {
		for i, val := range m[1:] {
			name := this.params[i]
			if req.URL.RawQuery != "" {
				req.URL.RawQuery += "&" + name +"="+ val
			} else {
				req.URL.RawQuery += name + "=" + val
			}
		}
	}
}

