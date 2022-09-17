# Terraform Provider for jq

This plugin brings the power of [jq](https://stedolan.github.io/jq/) to Terraform, allowing you to parse and extract complex data from JSON or even HCL objects.

## Usage

No initialiation is required, you can begin using the provider immediately.

```hcl
data "jq_query" "example" {
    data = "{\"a\": \"b\"}"
    query = ".a"
}

output "example" {
    value = data.jq_query.example.result
}
```

This will output:
```
Outputs:

example = "\"b\""
```

### HCL Compatibility

The jq operates on json formatted strings. Fortunately, terraform provides the `jsonencode()` and `jsondecode()` functions for easily converting back and forth between HCL and json strings.

The above example in pure HCL:

```hcl
data "jq_query" "example" {
    data = jsonencode({a = "b"})
    query = ".a"
}

output "example" {
    value = jsondecode(data.jq_query.example.result)
}
```

This will output:
```
Outputs:

example = "b"
```

### Multiple Results

Some jq queries will produce multiple elements. In this case, the result will have multiple lines, with each line containing an element as a JSON encoded string. Keep in mind this means that the result **cannot** be converted to HCL with `jsondecode()` as the string itself is not valid JSON.

```hcl
data "jq_query" "example" {
    data = jsonencode([{id:1},{id:2},{id:3}])
    query = ".[] | .id"
}

output "example" {
    value = data.jq_query.example.result
}
```

This will output:
```sh
Outputs:

example = <<EOT
1
2
3
EOT
```

Be sure to write your queries so they return a single element if you wish to convert them back to HCL with `jsondecode()`.
