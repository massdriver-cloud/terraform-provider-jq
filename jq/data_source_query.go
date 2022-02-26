package jq

import (
	"context"
	"strconv"
	"time"

	"github.com/savaki/jq"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceQuery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQueryRead,
		Schema: map[string]*schema.Schema{
			"data": {
				Type:        schema.TypeString,
				Description: "A JSON formatted string containing the data you would like to query. You can use jsonencode() to convert an HCL map or object into a queryable string.",
				Required:    true,
			},
			"query": {
				Type:        schema.TypeString,
				Description: "A jq query string.",
				Required:    true,
			},
			"result": {
				Type:        schema.TypeString,
				Description: "A JSON formatted string containing the result of the query",
				Computed:    true,
			},
		},
	}
}

func dataSourceQueryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	query, _ := jq.Parse(d.Get("query").(string))
	data := []byte(d.Get("data").(string))
	result, err := query.Apply(data)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("result", string(result)); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
