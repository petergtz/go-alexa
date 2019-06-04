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

type InteractionHistory interface {
	GetInteractionsByUser(userID string) []*Interaction
}

type InteractionLogger interface {
	Log(*Interaction)
}
