package restgo

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
)

// send json as response, you can give status code or default 200
func ServeJSON(ctx *fasthttp.RequestCtx, data interface{}, code int) {
	if code == 0 {
		code = http.StatusOK
	}

	content, err := json.Marshal(data)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.Write(content)
}

// send a string as response, you can give status code or default 200
func ServeTEXT(ctx *fasthttp.RequestCtx, data string, code int) {
	if code == 0 {
		code = http.StatusOK
	}
	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.Response.Header.Set("X-Content-Type-Options", "nosniff")
	ctx.SetStatusCode(code)
	ctx.Write([]byte(data))
}
