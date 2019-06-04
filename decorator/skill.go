package decorator

import (
	"time"

	"github.com/petergtz/go-alexa"
)

type InteractionLoggingSkill struct {
	delegate          alexa.Skill
	interactionLogger alexa.InteractionLogger
}

func ForSkillWithInteractionLogging(delegate alexa.Skill, interactionLogger alexa.InteractionLogger) *InteractionLoggingSkill {
	return &InteractionLoggingSkill{
		delegate:          delegate,
		interactionLogger: interactionLogger,
	}
}

func (s *InteractionLoggingSkill) ProcessRequest(requestEnv *alexa.RequestEnvelope) *alexa.ResponseEnvelope {
	s.interactionLogger.Log(&alexa.Interaction{
		RequestID:     requestEnv.Request.RequestID,
		RequestType:   requestEnv.Request.Type,
		UnixTimestamp: time.Now().Unix(),
		Timestamp:     time.Now(),
		Locale:        requestEnv.Request.Locale,
		UserID:        requestEnv.Session.User.UserID,
		SessionID:     requestEnv.Session.SessionID,
	})
	return s.delegate.ProcessRequest(requestEnv)
}
