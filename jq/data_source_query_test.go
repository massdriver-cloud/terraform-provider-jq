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
			result: "\"foo\"",
		},
		{
			data:   `{\"b\":1}`,
			query:  `.b`,
			result: "1",
		},
		{
			data:   `{\"c\":[\"foo\", \"bar\", \"baz\"]}`,
			query:  `.c.[1:2]`,
			result: "[\"bar\", \"baz\"]",
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
