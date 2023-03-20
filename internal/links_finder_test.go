package internal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinksFinder(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  struct {
			links []string
			err   error
		}
	}{
		{
			desc: "Receiving and html with links should return an slice of links",
			input: `<!DOCTYPE html>
		<html>
			<body>
				<h1>HTML Links</h1>
				<p><a href="https://www.test.com/">Test 1</a></p>
				<p><a href="https://www.test2.com/">Test 2</a></p>
			</body>
		</html>`,
			want: struct {
				links []string
				err   error
			}{
				links: []string{
					"https://www.test.com/",
					"https://www.test2.com/",
				},
			},
		},
		{
			desc:  "Receiving an empty string should return an empty slice of links",
			input: "",
			want: struct {
				links []string
				err   error
			}{
				links: []string{},
			},
		},
		{
			desc:  "Receiving something different to a html should return an empty slice of links",
			input: "this is not a html",
			want: struct {
				links []string
				err   error
			}{
				links: []string{},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte(tC.input))
			links, err := LinksFinder(buf)

			assert.Equal(t, tC.want.err, err)
			assert.Equal(t, tC.want.links, links)
		})
	}

}
