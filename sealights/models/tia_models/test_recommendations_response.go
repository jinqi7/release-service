package tia_models

import (
	"github.com/konflux-ci/release-service/__sealights__/models"
)

type TestRecommendationsExcludedTests struct {
	UniversalTestID string   `json:"universalTestId"`
	Name            string   `json:"name"`
	TestFramework   string   `json:"testFramework"`
	VariantIds      []string `json:"variantIds"`
}
type TestRecommendationsResponse struct {
	TestSelectionEnabled    bool                                `json:"testSelectionEnabled"`
	RecommendationSetStatus string                              `json:"recommendationSetStatus"`
	ExcludedTests           []*TestRecommendationsExcludedTests `json:"excludedTests"`
}

func NewTestRecommendationsResponse() *TestRecommendationsResponse {
	return &TestRecommendationsResponse{}
}

func (m *TestRecommendationsResponse) Validate() error {
	return nil
}

func (m *TestRecommendationsResponse) GetTestSelectionStatus(envConfigValue bool) string {
	if !envConfigValue {
		return "disabledByConfiguration"
	}
	if !m.TestSelectionEnabled {
		return "disabled"
	}
	switch m.RecommendationSetStatus {
	case "notReady":
		return "recommendationsTimeout"
	case "ready":
		return "recommendedTests"
	case "noHistory":
		return "disabled"
	case "error":
		return "error"
	case "wontBeReady":
		return "recommendationsTimeoutOnServer"
	default:
		return "error"
	}
}

var _ models.Model = (*TestRecommendationsResponse)(nil)
