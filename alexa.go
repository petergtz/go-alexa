package alexa

import "encoding/json"

type RequestEnvelope struct {
	Version string   `json:"version"`
	Session *Session `json:"session"`
	Request *Request `json:"request"`
	Context *Context `json:"context"`
}

func (requestEnv *RequestEnvelope) String() string {
	if requestEnv == nil {
		return "Not available."
	}
	r := *requestEnv
	r.Session.User.AccessToken = "<REDACTED>"
	buf, e := json.MarshalIndent(r, "", "  ")
	if e != nil {
		return "Error while marshalling request. Error: " + e.Error()
	}
	return string(buf)
}

// Session containes the session data from the Alexa request.
type Session struct {
	New         bool                   `json:"new"`
	SessionID   string                 `json:"sessionId"`
	Attributes  map[string]interface{} `json:"attributes"`
	User        User                   `json:"user"`
	Application Application            `json:"application"`
}

// Request contains the data in the request within the main request.
type Request struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Locale    string `json:"locale"`
	Timestamp string `json:"timestamp"`

	DialogState string `json:"dialogState"`

	// Only for LaunchRequest
	Target                     *LaunchRequestTarget   `json:"target"`
	Meta                       *LaunchRequestMetadata `json:"metadata"`
	Body                       interface{}            `json:"body"`
	Payload                    json.RawMessage        `json:"payload"`
	TargetUri                  string                 `json:"targetURI"`
	LaunchRequestType          string                 `json:"launchRequestType"`
	ShouldLinkResultBeReturned bool                   `json:"shouldLinkResultBeReturned"`

	// Only for IntentRequest
	Intent Intent `json:"intent"`

	// Only for SessionEndedRequest
	Reason string             `json:"reason"`
	Error  *SessionEndedError `json:"error"`

	// Only for Connections.Response
	Name   string                     `json:"name"`
	Status *ConnectionsResponseStatus `json:"status"`
	// Payload already defined for LaunchRequest
	Token string `json:"token"`
}

type LaunchRequestTarget struct {
	Path    string `json:"path"`
	Address string `json:"address"`
}

type LaunchRequestMetadata struct {
	Referrer string `json:"referrer"`
}

type SessionEndedError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Intent contains the data about the Alexa Intent requested.
type Intent struct {
	Name               string                `json:"name"`
	ConfirmationStatus string                `json:"confirmationStatus,omitempty"`
	Slots              map[string]IntentSlot `json:"slots"`
}

// IntentSlot contains the data for one Slot
type IntentSlot struct {
	Name               string                   `json:"name"`
	ConfirmationStatus string                   `json:"confirmationStatus,omitempty"`
	Value              string                   `json:"value"`
	Resolutions        *ResolutionsPerAuthority `json:"resolutions,omitempty"`
}

type ResolutionsPerAuthority struct {
	ResolutionsPerAuthority []Resolution `json:"resolutionsPerAuthority"`
}

type Resolution struct {
	Authority string            `json:"authority"`
	Status    map[string]string `json:"status"`
	Values    []Value           `json:"values"`
}

type Value struct {
	Value NameID `json:"value"`
}

type NameID struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type ConnectionsResponseStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ConnectionsResponsePayload struct {
	PurchaseResult string `json:"purchaseResult"`
	ProductID      string `json:"productId"`
	Message        string `json:"message"`
}

// ResponseEnvelope contains the Response and additional attributes.
type ResponseEnvelope struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          *Response              `json:"response"`
}

// Response contains the body of the response.
type Response struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	Directives       []interface{} `json:"directives,omitempty"`
	ShouldSessionEnd bool          `json:"shouldEndSession"`
}

// OutputSpeech contains the data the defines what Alexa should say to the user.
type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

// Card contains the data displayed to the user by the Alexa app.
type Card struct {
	Type    string `json:"type"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Text    string `json:"text,omitempty"`
	Image   *Image `json:"image,omitempty"`
}

// Image provides URL(s) to the image to display in resposne to the request.
type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// Reprompt contains data about whether Alexa should prompt the user for more data.
type Reprompt struct {
	OutputSpeech *OutputSpeech `json:"outputSpeech,omitempty"`
}

// AudioPlayerDirective contains device level instructions on how to handle the response.
type AudioPlayerDirective struct {
	Type         string     `json:"type"`
	PlayBehavior string     `json:"playBehavior,omitempty"`
	AudioItem    *AudioItem `json:"audioItem,omitempty"`
}

// AudioItem contains an audio Stream definition for playback.
type AudioItem struct {
	Stream Stream `json:"stream,omitempty"`
}

// Stream contains instructions on playing an audio stream.
type Stream struct {
	Token                string `json:"token"`
	URL                  string `json:"url"`
	OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
}

// DialogDirective contains directives for use in Dialog prompts.
type DialogDirective struct {
	Type          string  `json:"type"`
	SlotToElicit  string  `json:"slotToElicit,omitempty"`
	SlotToConfirm string  `json:"slotToConfirm,omitempty"`
	UpdatedIntent *Intent `json:"updatedIntent,omitempty"`
}

type APLRenderDocumentDirective struct {
	Type     string           `json:"type"`
	Token    string           `json:"token"`
	Document *json.RawMessage `json:"document"`
}

type ConnectionsSendRequestDirective struct {
	Type    string                        `json:"type"`
	Name    string                        `json:"name"`
	Payload ConnectionsSendRequestPayload `json:"payload"`
	Token   string                        `json:"token"`
}
type ConnectionsSendRequestPayload struct {
	InSkillProduct InSkillProduct `json:"InSkillProduct"`
}

type InSkillProduct struct {
	ProductID string `json:"productId"`
}

type Context struct {
	System *System `json:"System"`
}

type System struct {
	APIEndpoint    string       `json:"apiEndpoint"`
	APIAccessToken string       `json:"apiAccessToken"`
	Application    *Application `json:"application"`
	Device         *Device      `json:"device"`
	User           *User        `json:"user"`
}

type Application struct {
	ApplicationID string `json:"applicationId"`
}

type Device struct {
	DeviceID            string                 `json:"deviceId"`
	SupportedInterfaces map[string]interface{} `json:"supportedInterfaces"`
}

type User struct {
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken"`
}
