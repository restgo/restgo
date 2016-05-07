package router

import (
	"regexp"
	"strings"
)

func Path2Regexp(path string) (*regexp.Regexp, map[int]string){
	parts := strings.Split(path, "/")
	params := make(map[int]string)

	var j = 0
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			params[j] = strings.TrimPrefix(part, ":")
			parts[i] = "([^/]+)"
			j++
		}
	}

	path = strings.Join(parts, "/")

	return regexp.MustCompile(path), params
}
