package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Sched struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Event struct {
	UserID      string `json:"user_id"`
	EventID     string `json:"event_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Schedule    Sched  `json:"schedule"`
}

type ListEventsResponse struct {
	Events []Event `json:"events"`
}

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

func ListEvents(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("ListEvents")
	var (
		userID    = aws.String(request.RequestContext.Identity.CognitoIdentityID)
		tableName = aws.String(os.Getenv("EVENTS_TABLE_NAME"))
	)
	input := &dynamodb.ScanInput{
		TableName: tableName,
	}
	result, _ := ddb.Scan(input)

	var listEvents []Event
	for _, i := range result.Items {
		event := Event{}
		if err := dynamodbattribute.UnmarshalMap(i, &event); err != nil {
			fmt.Println("Failed to unmarshal")
			fmt.Println(err)
		}
		if event.UserID == *userID {
			listEvents = append(listEvents, event)
		}
	}
	body, _ := json.Marshal(&ListEventsResponse{
		Events: listEvents,
	})
	resp := events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
		Headers:    make(map[string]string),
	}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	return resp, nil
}

func main() {
	lambda.Start(ListEvents)
}
