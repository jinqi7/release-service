package tia_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type TestRecommendationsRequest struct {
	BuildSessionId string `json:"buildSessionId"`
	Stage          string `json:"stage"`
	TestGroupId    string `json:"testGroupId"`
}

func NewTestRecommendationsRequest() *TestRecommendationsRequest {
	return &TestRecommendationsRequest{}
}

func (m *TestRecommendationsRequest) Validate() error {
	if m.BuildSessionId == "" {
		return fmt.Errorf("TestRecommendationsRequest: buildSessionId is required")
	}
	if m.Stage == "" {
		return fmt.Errorf("TestRecommendationsRequest: stage is required")
	}
	return nil
}

func (m *TestRecommendationsRequest) SetBuildSessionId(buildSessionId string) *TestRecommendationsRequest {
	m.BuildSessionId = buildSessionId
	return m
}

func (m *TestRecommendationsRequest) SetStage(stage string) *TestRecommendationsRequest {
	m.Stage = stage
	return m
}

func (m *TestRecommendationsRequest) SetTestGroupId(testGroupId string) *TestRecommendationsRequest {
	m.TestGroupId = testGroupId
	return m
}

var _ models.Model = (*TestRecommendationsRequest)(nil)
