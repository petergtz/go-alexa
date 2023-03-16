package lambda

import (
	"context"

	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/aws/aws-lambda-go/lambda"
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
			"session-attributes", requestEnv.Session.Attributes,
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
			"session-attributes", result.SessionAttributes,
		)

		return result, nil
	})
}
