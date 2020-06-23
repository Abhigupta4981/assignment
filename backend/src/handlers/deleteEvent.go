package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var ddb *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		ddb = dynamodb.New(session)
	}
}

func deleteEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("DeleteEvent")
	var (
		userID    = aws.String(request.RequestContext.Identity.CognitoIdentityID)
		eventID   = aws.String(request.PathParameters["event_id"])
		tableName = aws.String(os.Getenv("EVENTS_TABLE_NAME"))
	)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: userID,
			},
			"event_id": {
				S: eventID,
			},
		},
		TableName: tableName,
	}
	_, err := ddb.DeleteItem(input)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		resp := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    make(map[string]string),
		}
		resp.Headers["Access-Control-Allow-Origin"] = "*"
		return resp, nil
	}
}

func main() {
	lambda.Start(deleteEvent)
}
