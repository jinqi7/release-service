package target

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type RequestType int

const (
	RequestTypeUndefined RequestType = iota
	RequestTypeAgentSendEvents
	RequestTypeGetExecutionId
	RequestTypeAgentSendFootprints
	RequestTypeAgentSendBuildMap
	RequestTypeAgentSendBuildMapEnd
	RequestTypeAgentSendTestEvents
	RequestTypeAgentSendExecutionStart
	RequestTypeAgentSendExecutionEnd
	RequestTypeAgentGetBuildSessionId
	RequestTypeAgentGetConfiguration
	RequestTypeAgentLogSubmission
	RequestTypeClockSync
	RequestTypeAgentGetTestRecommendations
	RequestTypeAgentGetBuildSessionIdByLabId
)

type TargetRequest interface {
	GetType() RequestType
	GetMetadata() Metadata
	GetData() models.Model
	Validate() error
}

func GetTypeRequestString(requestType RequestType) string {
	switch requestType {
	case RequestTypeAgentSendEvents:
		return "agent_send_event"
	case RequestTypeGetExecutionId:
		return "get_execution_id"
	case RequestTypeAgentSendFootprints:
		return "agent_send_footprints"
	case RequestTypeAgentSendBuildMap:
		return "agent_send_build_map"
	case RequestTypeAgentSendBuildMapEnd:
		return "agent_send_build_map_end"
	case RequestTypeAgentSendTestEvents:
		return "agent_send_test_events"
	case RequestTypeAgentSendExecutionStart:
		return "agent_send_execution_start"
	case RequestTypeAgentGetBuildSessionId:
		return "agent_get_build_session_id"
	case RequestTypeAgentSendExecutionEnd:
		return "agent_send_execution_end"
	case RequestTypeAgentGetConfiguration:
		return "agent_get_configuration"
	case RequestTypeAgentLogSubmission:
		return "agent_log_submission"
	case RequestTypeClockSync:
		return "clock_sync"
	case RequestTypeAgentGetTestRecommendations:
		return "agent_get_test_recommendations"
	case RequestTypeAgentGetBuildSessionIdByLabId:
		return "agent_get_build_session_id_by_lab_id"
	default:
		return "undefined"
	}
}

type Request struct {
	Type     RequestType
	Metadata Metadata
	Data     models.Model
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) GetType() RequestType {
	return r.Type
}

func (r *Request) GetMetadata() Metadata {
	return r.Metadata
}

func (r *Request) GetData() models.Model {
	return r.Data
}

func (r *Request) SetType(t RequestType) *Request {
	r.Type = t
	return r
}

func (r *Request) SetMetadata(m Metadata) *Request {
	r.Metadata = m
	return r
}

func (r *Request) SetData(data models.Model) *Request {
	r.Data = data
	return r
}

func (r *Request) SetMetadataKeyValue(key, value string) *Request {
	if r.Metadata == nil {
		r.Metadata = make(Metadata)
	}
	r.Metadata.Set(key, value)
	return r
}

func (r *Request) Validate() error {
	if r.Type == RequestTypeUndefined {
		return fmt.Errorf("Request: Type is undefined")
	}
	if r.Data == nil {
		return nil
	}
	if err := r.Data.Validate(); err != nil {
		return fmt.Errorf("Request: Data model is invalid: %s", err.Error())
	}
	return nil
}

var _ TargetRequest = (*Request)(nil)
