## Facet Frontend API

API consumed from Facet.Ninja Frontend applications.

## Getting started

Creating an executable zip file and uploading the file via terraform to AWS. AWS requires the executable location in the
file.

1. Run make:

```
cd /src/facet.ninja/api
make
```

2. Run [AWS Serverless Application Model (SAM)](https://aws.amazon.com/serverless/sam/) from the *root directory* of the
   project:

```
sam local start-api --env-vars env.json
```

3. Terraform is called a GH action.

## Run tests

`go test`

## Access JMeter results produced through CI

You can access JMeter tests via the S3
link: `https://cdn.facet.ninja/test/api/regression/PUT_THE_CI_NUMBER_HERE/index.html`,
i.e: https://cdn.facet.ninja/test/api/regression/88/index.html.

## Environment Variables

Environment variables are declared in both SAM (local development) and terraform modules (all the other environments).
These are the steps to declare an environment variable:

1. Declaring the variable under *Parameters*:

```
  MY_VARIABLE:
    Type: String
    Description: An example variable
    Default: This is the default example
```

2. Passing it into the lambda environment by reference:

```
Environment:
    Variables:
      MY_VARIABLE: !Ref MY_VARIABLE
```

3. Adding the actual value in `env.json`:

```
{
  "API": {
    "MY_VARIABLE": "This variable was loaded from env.json"
  }
}
```

4. Reading it in the application:

```
fmt.Println(os.Getenv("MY_VARIABLE"))
```

5. Declaring it in the terraform module [lambda.tf](./deploy/lambda.tf):

```
environment_variables       = {
    MY_VARIABLE = "This variable was loaded in Terraform"
  }
```

Read [here](https://github.com/aws/aws-sam-cli/issues/1163) for more about this pattern.

Everytime there is an environment variable change, please add its name and description
at [env-example-template.json](./env-example-template.json).

## System requirements

```
Go: 1.15.2
Java: 15
Docker desktop: 2.5
TF version: 0.13.5

```

## Terrafom Deployment

For deploying to the `test` environment, run `terraform apply` with current directory being `./deploy`. You will need to
have preconfigured SSH keys to access GH
repos. [Follow this guide](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/adding-a-new-ssh-key-to-your-github-account)
to generate SSH keys and upload them to GH. 