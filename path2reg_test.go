package restgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_path2reg(t *testing.T) {
	assert := assert.New(t)

	tests := [][]interface{}{
		{
			``,
			`\A\z`,
			true,
		},
		{
			`/`,
			`\A/\z`,
			true,
		},
		{
			`/foo/bar`,
			`\A/foo/bar\z`,
			true,
		},
		{
			`/foo/bar`,
			`\A/foo/bar\z`,
			true,
		},
		{
			`/:foo/bar/:baz`,
			`\A/(?P<foo>[^/#?]+)/bar/(?P<baz>[^/#?]+)\z`,
			false,
		},
		{
			`(/categories/:category_id)?/posts/:id`,
			`\A(/categories/(?P<category_id>[^/#?]+))?/posts/(?P<id>[^/#?]+)\z`,
			false,
		},
	}

	for _, pair := range tests {
		path := pair[0].(string)
		expectedRegexp := pair[1].(string)
		expectedIsStatic := pair[2].(bool)

		r, isStatic := path2Regexp(path, true)

		assert.Equal(expectedRegexp, r.String(), "path Regexp should match")
		assert.Equal(expectedIsStatic, isStatic, "path isStatic should match")
	}
}
