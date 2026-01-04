package lambda

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dustin/go-humanize"
	"github.com/petergtz/go-alexa"
)

func StartLambdaSkill(skill alexa.Skill, logger *zap.SugaredLogger) {
	invocationCount := 0
	lambda.Start(func(ctx context.Context, requestEnv alexa.RequestEnvelope) (alexa.ResponseEnvelope, error) {
		invocationCount++
		lc, _ := lambdacontext.FromContext(ctx)

		if requestEnv.Request == nil {

			logger.Infow("Keep-alive CloudWatch Request",
				"aws-request-id", lc.AwsRequestID,
				"function-invocation-count", invocationCount)

			return alexa.ResponseEnvelope{}, nil
		}
		logger.Infow("Alexa Request",
			"aws-request-id", lc.AwsRequestID,
			"alexa-request-id", requestEnv.Request.RequestID,
			"user-id", requestEnv.Session.User.UserID,
			"session-id", requestEnv.Session.SessionID,
			"locale", requestEnv.Request.Locale,
			"type", requestEnv.Request.Type,
			"intent", requestEnv.Request.Intent,
			"session-attributes", WithSnippets(requestEnv.Session.Attributes, 500),
			"function-invocation-count", invocationCount,
		)

		result := *skill.ProcessRequest(&requestEnv)

		logger.Infow("Alexa Response",
			"aws-request-id", lc.AwsRequestID,
			"alexa-request-id", requestEnv.Request.RequestID,
			"user-id", requestEnv.Session.User.UserID,
			"session-id", requestEnv.Session.SessionID,
			"locale", requestEnv.Request.Locale,
			"type", requestEnv.Request.Type,
			"intent", requestEnv.Request.Intent,
			"response", result.Response,
			"session-attributes", WithSnippets(result.SessionAttributes, 500),
		)

		return result, nil
	})
}

// Marshals every value into a JSON string and truncates it to maxLen characters, appending suffix if truncated
// if it's shorter then maxLen, it is left as is
func WithSnippets(s map[string]interface{}, maxLen int) map[string]interface{} {
	result := make(map[string]interface{}, len(s))
	for k, v := range s {
		b, e := json.Marshal(v)
		if e != nil {
			result[k] = "[error marshaling value: " + e.Error() + "]"
			continue
		}
		str := string(b)
		if len(str) > maxLen {
			result[k] = fmt.Sprintf("%s... (%s truncated)", str[:maxLen], humanize.Bytes(uint64(len(str)-maxLen)))
		} else {
			result[k] = v
		}
	}
	return result
}
