package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/konflux-ci/release-service/__sealights__/pkg/target/http/resty"
)

type DebugLog struct {
	Type      string      `json:"type"`
	URI       string      `json:"uri"`
	Method    string      `json:"method"`
	Timestamp time.Time   `json:"timestamp"`
	Headers   http.Header `json:"headers,omitempty"`
	Body      interface{} `json:"body,omitempty"`
	Error     string      `json:"error,omitempty"`
}

func (l *DebugLog) String() string {
	data, _ := json.MarshalIndent(l, "", " ")
	return string(data) + "\n"
}

func NewDebugLogForRequest(r *resty.RequestLog) *DebugLog {
	return &DebugLog{
		Type:      "request",
		URI:       r.RestyRequest.URL,
		Method:    r.RestyRequest.Method,
		Timestamp: time.Now(),
		Headers:   r.Header,
		Body:      r.RestyRequest.Body,
		Error:     "",
	}
}

func NewDebugLogForResponse(r *resty.ResponseLog) *DebugLog {
	return &DebugLog{
		Type:      "response",
		URI:       r.RestyResponse.Request.URL,
		Method:    r.RestyResponse.Request.Method,
		Timestamp: time.Now(),
		Headers:   r.Header,
		Body:      r.RestyResponse.Body(),
	}
}

func NewDebugLogForError(r *resty.Request, err error) *DebugLog {
	l := &DebugLog{
		Type:      "response",
		URI:       r.RawRequest.RequestURI,
		Method:    r.RawRequest.Method,
		Timestamp: time.Now(),
		Headers:   r.Header,
		Body:      r.Body,
		Error:     err.Error(),
	}
	return l
}
