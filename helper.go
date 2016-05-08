package grest

import (
	"net/http"
	"encoding/json"
	"fmt"
)

// send json as response, you can give status code or default 200
func ServeJSON(rw http.ResponseWriter, data interface{}, code int) {
	if code == 0 {
		code = http.StatusOK
	}

	content, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(content)
}

// send a string as response, you can give status code or default 200
func ServeTEXT(rw http.ResponseWriter, data string, code int) {
	if code == 0 {
		code = http.StatusOK
	}
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(code)
	fmt.Fprintln(rw, data)
}
