module "lambda" {
  source                      = "git@github.com:facets-io/terraform-module-aws-lambda-function.git?ref=0.0.4"
  function_name               = var.name
  handler                     = var.handler
  filename                    = var.filename
  runtime                     = var.runtime
  environment                 = var.environment
  environment_variables       = {
    COGNITO_JWKS_URL = "https://cognito-idp.us-west-2.amazonaws.com/us-west-2_oM4ne6cSf/.well-known/jwks.json"
  }
}
