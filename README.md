# Terraform Provider for jq

This plugin brings the power of [jq](https://stedolan.github.io/jq/) to Terraform, allowing you to parse and extract complex data from JSON or even HCL objects.

## Usage

No initialiation is required, you can begin using the provider immediately.

```
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

```
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
