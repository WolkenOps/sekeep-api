package main

import (
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

var (
	ssmClient ssmiface.SSMAPI
)

func init() {
	session := session.Must(session.NewSession())
	ssmClient = ssm.New(session)
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}
	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	return authResponse
}

func handler(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	tokenSlice := strings.Split(token, " ")
	var tokenString string
	if len(tokenSlice) > 1 {
		tokenString = tokenSlice[len(tokenSlice)-1]
	}

	parameter, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String("sekeep-token"),
	})

	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error getting token")
	}

	if tokenString != *parameter.Parameter.Value {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	return generatePolicy("user", "Allow", request.MethodArn), nil
}

func main() {
	lambda.Start(handler)
}
