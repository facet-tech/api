provider "aws" {
  allowed_account_ids = [var.account_id]
  region              = var.region
  version             = "2.61.0"
}

terraform {
  required_version = "= 0.13.0"

  backend "s3" {
    bucket  = "terraform-state.facets.ninja"
    key     = "facet.ninja/dev/api/terraform.tfstate"
    region  = "us-west-2"
    encrypt = true
  }
}
