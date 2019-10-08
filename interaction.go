package alexa

import "time"

type Interaction struct {
	RequestID     string                 `dynamodbav:"RequestID" json:"request_id"`
	RequestType   string                 `dynamodbav:"RequestType" json:"request_type"`
	UnixTimestamp int64                  `dynamodbav:"UnixTimestamp" json:"unix_timestamp"`
	Timestamp     time.Time              `dynamodbav:"Timestamp" json:"timestamp"`
	Locale        string                 `dynamodbav:"Locale" json:"locale"`
	UserID        string                 `dynamodbav:"UserID" json:"user_id"`
	SessionID     string                 `dynamodbav:"SessionID" json:"session_id"`
	Attributes    map[string]interface{} `dynamodbav:"Attributes" json:"attributes"`
}

func InteractionFrom(requestEnv *RequestEnvelope) *Interaction {
	timestamp := time.Now()
	i := &Interaction{
		RequestID:     requestEnv.Request.RequestID,
		RequestType:   requestEnv.Request.Type,
		UnixTimestamp: timestamp.Unix(),
		Timestamp:     timestamp,
		UserID:        requestEnv.Session.User.UserID,
		SessionID:     requestEnv.Session.SessionID,
		Locale:        requestEnv.Request.Locale,
	}
	if requestEnv.Request.Type == "IntentRequest" {
		return i.WithAttributes(map[string]interface{}{"Intent": requestEnv.Request.Intent.Name})
	}
	return i
}

func (i Interaction) WithAttributes(a map[string]interface{}) *Interaction {
	i.Attributes = a
	return &i
}

type InteractionHistory interface {
	GetInteractionsByUser(userID string, nwerThan time.Time) []*Interaction
}

type InteractionLogger interface {
	Log(*Interaction)
}
