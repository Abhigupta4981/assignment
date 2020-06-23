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
	uuid "github.com/satori/go.uuid"
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

func AddEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("AddEvent")
	var (
		userID    = aws.String(request.RequestContext.Identity.CognitoIdentityID)
		eventID   = uuid.Must(uuid.NewV4(), nil).String()
		tableName = aws.String(os.Getenv("EVENTS_TABLE_NAME"))
	)
	event := &Event{
		UserID:  *userID,
		EventID: eventID,
	}
	json.Unmarshal([]byte(request.Body), event)
	item, _ := dynamodbattribute.MarshalMap(event)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}
	if _, err := ddb.PutItem(input); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		body, _ := json.Marshal(event)
		resp := events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
			Headers:    make(map[string]string),
		}
		resp.Headers["Access-Control-Allow-Origin"] = "*"
		return resp, nil
	}
}
func main() {
	lambda.Start(AddEvent)
}
