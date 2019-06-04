package alexa

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Handler struct {
	Skill                 Skill
	Log                   *zap.SugaredLogger
	ExpectedApplicationID string
	SkipRequestValidation bool
}

type Skill interface {
	ProcessRequest(requestEnv *RequestEnvelope) *ResponseEnvelope
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	const timeLimit float64 = 150

	if !h.SkipRequestValidation && !IsValidAlexaRequest(w, req) {
		return
	}

	requestBody, e := ioutil.ReadAll(req.Body)
	if e != nil {
		h.Log.Errorw("Error while reading request body", "error", e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var alexaRequest RequestEnvelope
	e = json.Unmarshal(requestBody, &alexaRequest)
	if e != nil {
		h.Log.Errorw("Error while unmarshaling request body", "error", e)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if alexaRequest.Session == nil {
		h.Log.Infow("Session is empty", "error", e)
		http.Error(w, "Session is empty", http.StatusBadRequest)
		return
	}
	if alexaRequest.Session.Application.ApplicationID != h.ExpectedApplicationID {
		h.Log.Infof("ApplicationID does not match: %v", alexaRequest.Session.Application.ApplicationID)
		http.Error(w, "Invalid ApplicationID", http.StatusBadRequest)
		return
	}
	timestamp, e := time.Parse("2006-01-02T15:04:05Z", alexaRequest.Request.Timestamp)
	if e != nil {
		h.Log.Infof("Invalid timestamp. Timestamp: %v", alexaRequest.Request.Timestamp)
		http.Error(w, "Invalid Timestamp", http.StatusBadRequest)
		return
	}
	if math.Abs(time.Since(timestamp).Seconds()) > timeLimit {
		h.Log.Infow("Timestamp not within time limit.", "timestamp", alexaRequest.Request.Timestamp, "difference", math.Abs(time.Since(timestamp).Seconds()))
		http.Error(w, "Timestamp not within time limit", http.StatusBadRequest)
		return
	}

	output, e := json.Marshal(h.Skill.ProcessRequest(&alexaRequest))
	if e != nil {
		h.Log.Errorw("Error while marshalling response", "error", e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}
