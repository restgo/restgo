package restgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_match(t *testing.T) {
	assert := assert.New(t)

	tests := [][]interface{}{
		{
			`/`,  // router path
			true, // is end
			true, // static
			[][]interface{}{ // real url
				{`/`, true},
				{``, false},
			},
		},
		{
			`/`,
			false,
			true,
			[][]interface{}{ // real url
				{`/`, true},
				{`/test`, true},
				{``, false},
			},
		},
		{
			`/foo/bar`,
			true,
			true,
			[][]interface{}{
				{`/foo/bar`, true},
				{`/foo/bar/`, false},
				{`/foo/barr`, false},
				{`/fooo/bar`, false},
				{`/oo/bar`, false},
				{`/foo/`, false},
				{`/foo`, false},
			},
		},
		{
			`/:foo/bar/:baz`,
			true,
			false,
			[][]interface{}{
				{`/blog/bar/id`, true},
				{`/blog/bar/id/`, false},
				{`/blog/bar/ttttt`, true},
				{`/ttttt/bar/ttttt`, true},
				{`/blog/bar/`, false},
				{`/blog/id`, false},
				{`/blog/id`, false},
			},
		},
		{
			`/blog`, // router path
			false,   // is end
			true,    // static
			[][]interface{}{ // real url
				{`/blog`, true},
				{`/blog/article`, true},
				{`/blog/`, true},
				{`/bblog/`, false},
				{`/blogg`, true}, // for not end url, this can match
			},
		},
		{
			`/blog/:id`,
			false,
			false,
			[][]interface{}{
				{`/blog/id`, true},
				{`/blog/asd/`, true},
				{`/blog/asd/asd`, true},
				{`/blog/asd/asd/asd`, true},
				{`/blog/`, false},
				{`/bblog/`, false},
				{`/blogg`, false},
			},
		},
		{
			`/:blog`,
			false,
			false,
			[][]interface{}{
				{`/b1/id`, true},
				{`/b2/asd/`, true},
				{`/b3/asd/asd`, true},
				{`/b4/asd/asd/asd`, true},
				{`/b5/`, true},
				{`/b6`, true},
			},
		},
	}

	for _, pair := range tests {
		routerPath, _ := pair[0].(string)
		isEnd, _ := pair[1].(bool)
		l := newLayer(routerPath, nil, isEnd)

		// isStatic
		assert.Equal(pair[2], l.isStatic, "isStatic should be", pair[2])

		// match
		urls, _ := pair[3].([][]interface{})
		for _, url := range urls {
			testUrl, _ := url[0].(string)
			expectedMatch, _ := url[1].(bool)
			_, match := l.match(testUrl)
			assert.Equal(expectedMatch, match, testUrl, "should match with ", l.pathRegexp.String())

		}
	}

}
