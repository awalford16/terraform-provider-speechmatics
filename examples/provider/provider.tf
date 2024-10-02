terraform {
  required_providers {
    speechmatics = {
      source = "hashicorp.com/edu/speechmatics"
    }
  }
}

provider "speechmatics" {
  endpoint = "asr.api.speechmatics.com"
}

data "speechmatics_example" "example" {}
