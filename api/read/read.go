package main

import (
	"github.com/WolkenOps/sekeep-api/internal/manager"
	"github.com/WolkenOps/sekeep-api/internal/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	password := model.Password{}

	password.Name = request.QueryStringParameters["name"]
	value, err := manager.Read(password)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: err.StatusCode,
			Body:       err.Message,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       value,
	}, nil
}

func main() {
	lambda.Start(handler)
}
