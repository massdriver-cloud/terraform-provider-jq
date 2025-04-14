terraform {
  required_providers {
    jq = {
      version = "0.2.1"
      source  = "massdriver-cloud/jq"
    }
  }
}

provider "jq" {}

data "jq_query" "example" {
  data = jsonencode({a = "b"})
  query = ".a"
}

output "example" {
  value = jsondecode(data.jq_query.example.result)
}
