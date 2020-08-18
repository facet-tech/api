module "lambda" {
  source                      = "git@github.com:facets-io/terraform-module-aws-lambda-function.git?ref=0.0.4"
  function_name               = var.name
  handler                     = var.handler
  filename                    = var.filename
  runtime                     = var.runtime
  environment                 = var.environment
  environment_variables       = {}
}
