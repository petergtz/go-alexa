package dynamodb

import (
	"fmt"
	"time"

	"github.com/petergtz/go-alexa"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"
)

type RequestLogger struct {
	dynamo          *dynamodb.DynamoDB
	logger          *zap.SugaredLogger
	tableName       string
	interactionChan chan *alexa.Interaction
}

func NewInteractionLogger(dynamoClient *dynamodb.DynamoDB, logger *zap.SugaredLogger, tableName string) *RequestLogger {
	p := &RequestLogger{
		dynamo:          dynamoClient,
		logger:          logger,
		tableName:       tableName,
		interactionChan: make(chan *alexa.Interaction, 100),
	}

	// Logging asynchronously to minimize latency when handling Alexa requests.
	// Since we get very few requests at this point, it's perfectly fine to do this with a single worker.
	go func() {
		for interaction := range p.interactionChan {
			p.doLog(interaction)
		}
	}()

	return p
}

func (p *RequestLogger) Log(interaction *alexa.Interaction) {
	p.interactionChan <- interaction
}

func (p *RequestLogger) doLog(interaction *alexa.Interaction) {
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

func (p *RequestLogger) GetInteractionsByUser(userID string, newerThan time.Time) []*alexa.Interaction {
	result := make([]*alexa.Interaction, 0)
	e := p.dynamo.QueryPages(&dynamodb.QueryInput{
		IndexName:              aws.String("UserID-UnixTimestamp-index"),
		KeyConditionExpression: aws.String("UserID = :userID and UnixTimestamp > :timestamp"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID":    &dynamodb.AttributeValue{S: &userID},
			":timestamp": &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%v", newerThan.Unix()))},
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
