package target

import (
	"github.com/konflux-ci/release-service/__sealights__/models"
)

type ResponseType int

const (
	ResponseTypeUndefined ResponseType = iota
	ResponseTypeGeneral
	ResponseTypeGetExecution
	ResponseTypesGetBuildSessionId
	ResponseTypesGetAgentConfiguration
	ResponseTypesClockSync
	ResponseTypesTestRecommendations
)

type TargetResponse interface {
	GetType() ResponseType
	GetMetadata() Metadata
	GetData() models.Model
	GetError() error
	Validate() error
}

type Response struct {
	Type     ResponseType
	Metadata Metadata
	Data     models.Model
	Error    error
}

func (r *Response) GetType() ResponseType {
	return r.Type
}

func (r *Response) GetMetadata() Metadata {
	return r.Metadata
}

func (r *Response) GetData() models.Model {
	return r.Data
}

func (r *Response) GetError() error {
	return r.Error
}

func (r *Response) Validate() error {
	return nil
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) SetType(t ResponseType) *Response {
	r.Type = t
	return r
}

func (r *Response) SetMetadata(value Metadata) *Response {
	r.Metadata = value
	return r
}

func (r *Response) SetMetadataKeyValue(key, value string) *Response {
	if r.Metadata == nil {
		r.Metadata = make(Metadata)
	}
	r.Metadata.Set(key, value)
	return r
}

func (r *Response) SetData(value models.Model) *Response {
	r.Data = value
	return r
}

func (r *Response) SetError(err error) *Response {
	r.Error = err
	return r
}

var _ TargetResponse = (*Response)(nil)
