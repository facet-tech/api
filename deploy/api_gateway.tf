module "api_gateway" {
  source                 = "git@github.com:facets-io/terraform-module-aws-api-gateway.git?ref=0.0.2"
  name                   = "${var.name}"
  environment            = var.environment
  description            = var.description
  versioned_directory    = path.cwd
  deploy_test_stage      = var.deploy_test_stage
  deploy_live_stage      = var.deploy_live_stage
  create_custom_domain   = true
  route53_record_zone_id = var.route53_record_zone_id
  route53_record_name    = var.route53_record_name

  endpoints = [
    {
      path    = "/'{proxy+}'"
      methods = [
        {
          method              = "ANY"
          integration_request = {
            uri                  = module.lambda.lambda_versioned_invoke_arn
          }
        }
      ]
    }
  ]
}