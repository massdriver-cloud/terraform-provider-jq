package jq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQueryBasic(t *testing.T) {
	type testData struct {
		data   string
		query  string
		result string
	}
	tests := []testData{
		{
			data:   `{\"a\":\"foo\"}`,
			query:  `.a`,
			result: `"foo"`,
		},
		{
			data:   `{\"b\":1}`,
			query:  `.b`,
			result: "1",
		},
		{
			data:   `{\"a\":1,\"b\":2}`,
			query:  `.a += 1 | .b *= 2`,
			result: "{\"a\":2,\"b\":4}",
		},
		{
			data:  `[{\"id\":1},{\"id\":2},{\"id\":3}]`,
			query: `.[] | .id`,
			result: `1
2
3`,
		},
		{
			data:   `[{\"id\":1},{\"id\":2},{\"id\":3}]`,
			query:  `[.[].id]`,
			result: `[1,2,3]`,
		},
		{
			data:   `{\"a\":{\"z\":{\"foo\":\"bar\"}},\"b\":{\"x\":[9]},\"c\":{\"z\":{\"hello\":\"world\"}}}`,
			query:  `[ .[].z | select( . != null )]`,
			result: `[{"foo":"bar"},{"hello":"world"}]`,
		},
		{
			data:   `{\"c\":[\"foo\", \"bar\", \"baz\"]}`,
			query:  `.c.[1:3]`,
			result: "[\"bar\",\"baz\"]",
		},
	}

	for _, tc := range tests {
		resource.Test(t, resource.TestCase{
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccCheckJqQueryConfig(tc.data, tc.query),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.jq_query.new", "result", tc.result),
					),
				},
			},
		})
	}
}

func testAccCheckJqQueryConfig(data, query string) string {
	return fmt.Sprintf(`
	data "jq_query" "new" {
		data = "%s"
		query = "%s"
	}
	`, data, query)
}
