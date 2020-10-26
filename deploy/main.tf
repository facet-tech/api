provider "aws" {
  allowed_account_ids = [var.account_id]
  region              = var.region
  version             = "2.61.0"
}

terraform {
  required_version = "= 0.13.5"

  backend "s3" {
    bucket  = "terraform-state.facets.ninja"
    key     = "facet.ninja/dev/api/terraform.tfstate"
    region  = "us-west-2"
    encrypt = true
  }
}

provider "external" {
  version = "2.0.0"
}

provider "null" {
  version = "3.0.0"
}