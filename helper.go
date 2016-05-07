package router

import (
	"net/http"
	"encoding/json"
)

func ServeJSON(rw http.ResponseWriter, data interface{}) {
	content, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(content)
}
