account_id                   = "935571265336"
aws_lambda_qualifier_version = "$LATEST"
deploy_live_stage            = false
deploy_test_stage            = true
description                  = "Facet API"
environment                  = "prod"
filename                     = "../build/main.zip"
handler                      = "./main"
name                         = "api"
region                       = "us-west-2"
route53_record_name          = "api"
route53_record_zone_id       = "Z05672452AKEG6MP6GI8Y"
runtime                      = "go1.x"