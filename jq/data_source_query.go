package jq

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itchyny/gojq"
)

func dataSourceQuery() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to execute a jq query against a JSON formatted string.",
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
				Description: "A JSON formatted string containing the result of the query.",
				Computed:    true,
			},
		},
	}
}

func dataSourceQueryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	result := ""
	var data interface{}
	var iter gojq.Iter
	query, err := gojq.Parse(d.Get("query").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = json.Unmarshal([]byte(d.Get("data").(string)), &data)
	if err != nil {
		return diag.FromErr(err)
	}

	// this the itchyny library requires the input to be a map[string]interface{} or []interface
	dataMap, ok := data.(map[string]interface{})
	if ok {
		iter = query.Run(dataMap)
	} else {
		dataArray, ok := data.([]interface{})
		if ok {
			iter = query.Run(dataArray)
		} else {
			return diag.FromErr(errors.New(`unable to process "data" argument: must be a json object or array`))
		}
	}

	count := 0
	// iterate through results and compile string
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return diag.FromErr(err)
		}
		resultBytes, err := json.Marshal(v)
		if err != nil {
			return diag.FromErr(err)
		}
		if count > 0 {
			result += "\n"
		}
		result += string(resultBytes)
		count++
	}

	if err := d.Set("result", string(result)); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
