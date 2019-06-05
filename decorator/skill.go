package decorator

import (
	"github.com/petergtz/go-alexa"
)

type InteractionLoggingSkill struct {
	delegate          alexa.Skill
	interactionLogger alexa.InteractionLogger
	shouldLog         func(requestEnv *alexa.RequestEnvelope) bool
}

func ForSkillWithInteractionLogging(
	delegate alexa.Skill,
	interactionLogger alexa.InteractionLogger,
	shouldLogFunc func(requestEnv *alexa.RequestEnvelope) bool) *InteractionLoggingSkill {
	return &InteractionLoggingSkill{
		delegate:          delegate,
		interactionLogger: interactionLogger,
		shouldLog:         shouldLogFunc,
	}
}

func (s *InteractionLoggingSkill) ProcessRequest(requestEnv *alexa.RequestEnvelope) *alexa.ResponseEnvelope {
	if s.shouldLog != nil && s.shouldLog(requestEnv) {
		s.interactionLogger.Log(alexa.InteractionFrom(requestEnv))
	}
	return s.delegate.ProcessRequest(requestEnv)
}
