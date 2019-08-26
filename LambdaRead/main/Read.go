package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
)

const BUCKET = "beier-wayne-test"
const region = "us-east-1"

// to connect the database and to use get method to get the device with id from database
func getDevices(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var filePath = request.PathParameters["file_path"]
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error," + err.Error(),
			StatusCode: 500,
		}, nil
	}

	// Create s3 connection
	svc := s3.New(sess)
	param := &s3.GetObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(filePath),
	}

	result, err := svc.GetObject(param)
	if err != nil {
			if err != nil {
				return events.APIGatewayProxyResponse{
					Body:       "Internal Server Error," + err.Error(),
					StatusCode: 500,
				}, nil
			}
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
			if err != nil {
				return events.APIGatewayProxyResponse{
					Body:       "Can't read files," + err.Error(),
					StatusCode: 500,
				}, nil
			}
	}
    return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(getDevices)
}
