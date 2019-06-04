package dynamodb

import (
	"github.com/petergtz/go-alexa"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"
)

type RequestLogger struct {
	dynamo    *dynamodb.DynamoDB
	logger    *zap.SugaredLogger
	tableName string
}

func NewInteractionLogger(accessKeyID, secretAccessKey string, logger *zap.SugaredLogger, tableName string) *RequestLogger {

	dynamoClient := dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})))

	return &RequestLogger{
		dynamo:    dynamoClient,
		logger:    logger,
		tableName: tableName,
	}
}

func (p *RequestLogger) Log(interaction *alexa.Interaction) {
	input, e := dynamodbattribute.MarshalMap(interaction)
	if e != nil {
		p.logger.Errorw("Could not marshal entry", "error", e)
		return
	}
	_, e = p.dynamo.PutItem(&dynamodb.PutItemInput{
		Item:      input,
		TableName: &p.tableName,
	})
	if e != nil {
		p.logger.Errorw("Could not log requests", "error", e)
		return
	}
}

func (p *RequestLogger) GetInteractionsByUser(userID string) []*alexa.Interaction {
	result := make([]*alexa.Interaction, 0)
	e := p.dynamo.QueryPages(&dynamodb.QueryInput{
		IndexName:              aws.String("UserID-UnixTimestamp-index"),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": &dynamodb.AttributeValue{S: &userID},
		},
		TableName: &p.tableName,
	},
		func(output *dynamodb.QueryOutput, lastPage bool) bool {
			for _, item := range output.Items {
				var interaction alexa.Interaction
				e := dynamodbattribute.UnmarshalMap(item, &interaction)
				if e != nil {
					continue
				}
				result = append(result, &interaction)
			}
			return true
		})
	if e != nil {
		p.logger.Errorw("Could query table", "error", e)
		return nil
	}
	return result
}
