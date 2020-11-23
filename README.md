## facet.ninja API

API for the facet.ninja extension.

## Getting started

1. Run make:

```
cd /src/facet.ninja/api
make
```

2. Run [AWS Serverless Application Model (SAM)](https://aws.amazon.com/serverless/sam/) from the *root directory* of the project: 

```
sam local start-api
```

## Run tests

`go test`

## Access JMeter results produced through CI

You can access JMeter tests via the S3 link: `https://cdn.facet.ninja/test/api/regression/PUT_THE_CI_NUMBER_HERE/index.html`, i.e: https://cdn.facet.ninja/test/api/regression/88/index.html. 