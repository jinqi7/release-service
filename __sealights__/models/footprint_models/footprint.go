package footprint_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
	"github.com/konflux-ci/release-service/__sealights__/models/agent_models"
)

type FootprintMeta struct {
	AgentMetadata *agent_models.AgentMetadata `json:"agentMetadata,omitempty"`
	AgentConfig   *agent_models.AgentConfig   `json:"agentConfig,omitempty"`
	AgentId       string                      `json:"agentId"`
	LabId         string                      `json:"labId"`
	Intervals     struct {
		TimedFootprintsCollectionIntervalSeconds int64 `json:"timedFootprintsCollectionIntervalSeconds"`
	} `json:"intervals"`
}
type Footprint struct {
	FormatVersion         string                `json:"formatVersion"`
	Methods               []string              `json:"methods"`
	Executions            []*FootprintExecution `json:"executions"`
	Meta                  FootprintMeta         `json:"meta"`
	modelMap              map[string]int
	executionId           string
	allowEmptyExecutionId bool
}

func NewFootprint() *Footprint {
	fm := &Footprint{
		FormatVersion: "6.0",
		Methods:       nil,
		Executions:    nil,
		Meta: FootprintMeta{
			AgentMetadata: nil,
			AgentConfig:   nil,
			AgentId:       "",
			LabId:         "",
			Intervals: struct {
				TimedFootprintsCollectionIntervalSeconds int64 `json:"timedFootprintsCollectionIntervalSeconds"`
			}{
				TimedFootprintsCollectionIntervalSeconds: 0,
			},
		},
		executionId:           "",
		modelMap:              make(map[string]int),
		allowEmptyExecutionId: false,
	}
	return fm
}

func (f *Footprint) SetAgentId(agentId string) *Footprint {
	f.Meta.AgentId = agentId
	return f
}

func (f *Footprint) SetLabId(labId string) *Footprint {
	f.Meta.LabId = labId
	return f
}

func (f *Footprint) SetTimedFootprintsCollectionIntervalSeconds(interval int64) *Footprint {
	f.Meta.Intervals.TimedFootprintsCollectionIntervalSeconds = interval
	return f
}

func (f *Footprint) SetExecutionId(executionId string) *Footprint {
	f.executionId = executionId
	return f
}

func (f *Footprint) SetAllowEmptyExecutionId(allowEmptyExecutionId bool) *Footprint {
	f.allowEmptyExecutionId = allowEmptyExecutionId
	return f
}

func (f *Footprint) SetAgentMetadata(value *agent_models.AgentMetadata) *Footprint {
	f.Meta.AgentMetadata = value
	return f
}

func (f *Footprint) SetAgentConfig(value *agent_models.AgentConfig) *Footprint {
	f.Meta.AgentConfig = value
	return f
}

func (f *Footprint) AddAgentMetadata(isAdd bool) *Footprint {
	if isAdd {
		f.Meta.AgentMetadata = agent_models.NewAgentMetadata()
	} else {
		f.Meta.AgentMetadata = nil
	}
	return f
}

func (f *Footprint) Validate() error {
	if f.FormatVersion == "" {
		return fmt.Errorf("Footprint: formatVersion is empty")
	}
	if len(f.Executions) == 0 {
		return fmt.Errorf("Footprint: executions is empty")
	}
	for _, execution := range f.Executions {
		execution.SetAllowEmptyExecution(f.allowEmptyExecutionId)
		if err := execution.Validate(); err != nil {
			return err
		}
	}
	if f.Methods == nil {
		return fmt.Errorf("Footprint: methods is empty")
	}
	if f.Meta.AgentId == "" {
		return fmt.Errorf("Footprint: agentId is empty")
	}
	if f.Meta.LabId == "" {
		return fmt.Errorf("Footprint: labId is empty")
	}
	if f.Meta.Intervals.TimedFootprintsCollectionIntervalSeconds <= 0 {
		return fmt.Errorf("Footprint: timedFootprintsCollectionIntervalSeconds is empty")
	}

	return nil
}

func (f *Footprint) AddExecution(execution *FootprintExecution) *Footprint {
	f.Executions = append(f.Executions, execution)
	return f
}

func (f *Footprint) AddMethod(method string) *Footprint {
	if _, ok := f.modelMap[method]; !ok {
		f.Methods = append(f.Methods, method)
		f.modelMap[method] = len(f.Methods) - 1
	}
	return f
}

func (f *Footprint) AddFootprint(activeTests []string, methodName string, startTimestamp, endTimestamp int64, isInit bool) *Footprint {
	if len(f.Executions) == 0 {
		f.Executions = append(f.Executions, NewFootprintExecution().SetExecutionId(f.executionId))
	}
	f.AddMethod(methodName)
	fe := f.Executions[len(f.Executions)-1]
	if len(activeTests) == 0 {
		hit := fe.FindHit("", startTimestamp, endTimestamp)
		if hit == nil {
			fe.AddHit(NewFootprintExecutionHit().
				SetStart(startTimestamp).SetEnd(endTimestamp).
				SetIsInitFootprints(isInit).
				SetMethods(f.modelMap[methodName]))
		} else {
			hit.SetMethods(f.modelMap[methodName])
		}
	} else {
		for _, test := range activeTests {
			hit := fe.FindHit(test, startTimestamp, endTimestamp)
			if hit == nil {
				testName := new(string)
				*testName = test
				fe.AddHit(NewFootprintExecutionHit().
					SetTestName(testName).
					SetStart(startTimestamp).SetEnd(endTimestamp).
					SetIsInitFootprints(isInit).
					SetMethods(f.modelMap[methodName]))

			} else {
				hit.SetMethods(f.modelMap[methodName])
			}
		}
	}
	return f
}

func (f *Footprint) updateMethodMap() *Footprint {
	f.modelMap = make(map[string]int)
	for i, method := range f.Methods {
		f.modelMap[method] = i
	}
	return f
}

func (f *Footprint) updateTimeInterval() *Footprint {
	f.Executions[0].SortHits()
	start, end := f.Executions[0].GetTimestampSpan()
	seconds := roundTimeToSeconds(start, end)
	return f.SetTimedFootprintsCollectionIntervalSeconds(seconds)
}

func (f *Footprint) Complete(aggType FootprintAggregationMode) *Footprint {
	switch aggType {
	case FootprintAggregationModeAllHits, FootprintAggregationModeCoverage:
	case FootprintAggregationModeAggregateTimeSpan:
		f.Executions[0].AggregateTimestamps()
	default:

	}

	return f.updateTimeInterval()
}

func (f *Footprint) Stats() (int, int, int64) {
	return len(f.Methods), len(f.Executions[0].Hits), f.Meta.Intervals.TimedFootprintsCollectionIntervalSeconds
}

func roundTimeToSeconds(start, end int64) int64 {
	return ((end - start) / 1e3) + 1
}

var _ models.Model = (*Footprint)(nil)
