# AWS_Exercise

## Description

  simple AWS Exercise

------
## Dependency

| package             |            method                       |
|---------------------|-----------------------------------------|
| go                  | sudo yum install -y golang              |
| AWS                 | github.com/aws/aws-sdk-go/aws           |
| MUX                 | github.com/gorilla/mux                  |

------

## <1> Deploy
### 1. Preparation
1. IAM Role for lambda permision
2. Get access keys and export them
3. set the bitbucket name in three files or as env(to do).

------
### 2. deploy and build
1. make
2. serverless deploy

------
## <2> Run
go build && ./rest-api to start the server
Test on the AWS Lambda or use pacman to query the api

------

## <3> Clean up
	serverless remove 