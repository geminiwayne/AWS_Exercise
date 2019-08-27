package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/satori/go.uuid"
)

type Client struct {
	Name       string `json:"Name"`
   Model      string `json:"Model"`
}

const BUCKET = "xxxx"
const region = "us-east-1"

func CreateModel(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error," + err.Error(),
			StatusCode: 500,
		}, nil
	}
	client := Client{}
	err = json.Unmarshal([]byte(request.Body), &client)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Bad Request," + err.Error(),
			StatusCode: 400,
		}, nil
	}
	uid, err := uuid.NewV4()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Can't create file name," + err.Error(),
			StatusCode: 500,
		}, nil
	}else {
		var fileName = uid.String()
		err = WriteFiles(client, fileName+".json", sess)
		if err != nil{
			return events.APIGatewayProxyResponse{
				Body:       "Can't write files to s3," + err.Error(),
				StatusCode: 500,
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "Client Information Create Successfully",
		StatusCode: 200,
	}, nil
}

// write files to s3
func WriteFiles(client Client, path string, svc *session.Session) error {
	var result []byte
	data, dataErr := json.Marshal(client)
	if dataErr != nil {
			return dataErr
		}
	result = append(result, data...)

	_, err := s3manager.NewUploader(svc).Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(path),
		Body:   bytes.NewReader(result),
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(CreateModel)
}
