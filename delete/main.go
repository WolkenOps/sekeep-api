package main

import (
	"encoding/json"

	"github.com/WolkenOps/sekeep-api/lib/manager"
	"github.com/WolkenOps/sekeep-api/lib/model"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	password := model.Password{}
	json.Unmarshal([]byte(request.Body), &password)

	err := manager.Delete(password)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: err.StatusCode,
			Body:       err.Message,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       password.Name,
	}, nil
}

func main() {
	lambda.Start(handler)
}
