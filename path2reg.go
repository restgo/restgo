package restgo

import (
	"fmt"
	"regexp"
)

// refer to https://github.com/pilu/traffic/blob/master/utils.go
func path2Regexp(path string, end bool) (*regexp.Regexp, bool) {
	var re *regexp.Regexp
	var isStatic bool

	regexpString := path

	isStaticRegexp := regexp.MustCompile(`[\(\)\?\<\>:]`)
	if !isStaticRegexp.MatchString(path) {
		isStatic = true
	}

	// Dots
	re = regexp.MustCompile(`([^\\])\.`)
	regexpString = re.ReplaceAllStringFunc(regexpString, func(m string) string {
		return fmt.Sprintf(`%s\.`, string(m[0]))
	})

	// Wildcard names
	re = regexp.MustCompile(`:[^/#?()\.\\]+\*`)
	regexpString = re.ReplaceAllStringFunc(regexpString, func(m string) string {
		return fmt.Sprintf("(?P<%s>.+)", m[1:len(m)-1])
	})

	re = regexp.MustCompile(`:[^/#?()\.\\]+`)
	regexpString = re.ReplaceAllStringFunc(regexpString, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)])
	})

	var str string
	if end {
		str = fmt.Sprintf(`\A%s\z`, regexpString)
	} else {
		str = fmt.Sprintf(`\A%s/?`, regexpString)
	}

	return regexp.MustCompile(str), isStatic
}
