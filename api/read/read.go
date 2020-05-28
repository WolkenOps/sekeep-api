package main

import (
	"encoding/json"
	"net/url"

	"github.com/WolkenOps/sekeep-api/internal/manager"
	"github.com/WolkenOps/sekeep-api/internal/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	password := model.Password{}

	if value, ok := request.PathParameters["name"]; ok {
		password.Name, _ = url.QueryUnescape(value)
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

	password.Name = request.QueryStringParameters["filter"]
	value, err := manager.List(password)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: err.StatusCode,
			Body:       err.Message,
		}, nil
	}

	passwords, _ := json.Marshal(value)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(passwords),
	}, nil
}

func main() {
	lambda.Start(handler)
}
