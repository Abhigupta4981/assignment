package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"gopkg.in/go-playground/validator.v9"
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

type Sched struct {
	StartTime *string `json:"start_time,omitempty"`
	EndTime   *string `json:"end_time,omitempty"`
}

type Event struct {
	UserID      *string `json:"user_id,omitempty"`
	EventID     *string `json:"event_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	Schedule    Sched   `json:"schedule,omitempty"`
}

func completeEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("CompleteEvent")
	var (
		userID    = aws.String(request.RequestContext.Identity.CognitoIdentityID)
		eventID   = aws.String(request.PathParameters["event_id"])
		tableName = aws.String(os.Getenv("EVENTS_TABLE_NAME"))
	)

	event := &Event{}
	json.Unmarshal([]byte(request.Body), event)

	var validate *validator.Validate
	validate = validator.New()
	er := validate.Struct(event)

	if er != nil {
		return events.APIGatewayProxyResponse{
			Body:       er.Error(),
			StatusCode: 400,
		}, nil
	}

	update := expression.UpdateBuilder{}
	u := reflect.ValueOf(event).Elem()
	t := u.Type()

	for i := 0; i < u.NumField(); i++ {
		f := u.Field(i)
		if !reflect.DeepEqual(f.Interface(), reflect.Zero(f.Type()).Interface()) {
			jsonFieldName := t.Field(i).Name
			if jsonTag := t.Field(i).Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
					jsonFieldName = jsonTag[:commaIdx]
				}
			}
			update = update.Set(expression.Name(jsonFieldName), expression.Value(f.Interface()))
		}
	}
	builder := expression.NewBuilder().WithUpdate(update)
	expression, erro := builder.Build()

	if erro != nil {
		return events.APIGatewayProxyResponse{
			Body:       erro.Error(),
			StatusCode: 400,
		}, nil
	}
	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: userID,
			},
			"event_id": {
				S: eventID,
			},
		},
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		UpdateExpression:          expression.Update(),
		ReturnValues:              aws.String("UPDATED_NEW"),
		TableName:                 tableName,
	}
	_, err := ddb.UpdateItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		resp := events.APIGatewayProxyResponse{
			Body:       request.Body,
			StatusCode: 200,
			Headers:    make(map[string]string),
		}
		resp.Headers["Access-Control-Allow-Origin"] = "*"
		return resp, nil
	}
}
func main() {
	lambda.Start(completeEvent)
}
