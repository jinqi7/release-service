package footprint_models

import (
	"fmt"
	"sort"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type FootprintExecution struct {
	ExecutionId         string                   `json:"executionId"`
	Hits                []*FootprintExecutionHit `json:"hits"`
	hitsMap             map[string]*FootprintExecutionHit
	allowEmptyExecution bool
}

func NewFootprintExecution() *FootprintExecution {
	return &FootprintExecution{
		ExecutionId: "",
		Hits:        nil,
		hitsMap:     make(map[string]*FootprintExecutionHit),
	}
}

func (f *FootprintExecution) SetExecutionId(executionId string) *FootprintExecution {
	f.ExecutionId = executionId
	return f
}

func (f *FootprintExecution) SetAllowEmptyExecution(allowEmptyExecution bool) *FootprintExecution {
	f.allowEmptyExecution = allowEmptyExecution
	return f
}

func (f *FootprintExecution) Validate() error {
	if !f.allowEmptyExecution && f.ExecutionId == "" {
		return fmt.Errorf("FootprintExecution: executionId is empty")
	}
	if len(f.Hits) == 0 {
		return fmt.Errorf("FootprintExecution: hits is empty")
	}
	for _, hit := range f.Hits {
		if err := hit.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (f *FootprintExecution) FindHit(testName string, startTimestamp, endTimestamp int64) *FootprintExecutionHit {
	key := fmt.Sprintf("%s%d%d", testName, startTimestamp, endTimestamp)
	return f.hitsMap[key]
}

func (f *FootprintExecution) AddHit(hits ...*FootprintExecutionHit) *FootprintExecution {
	for _, hit := range hits {
		f.Hits = append(f.Hits, hit)
		f.hitsMap[hit.Key()] = hit
	}
	return f
}

func (f *FootprintExecution) SortHits() *FootprintExecution {
	sort.Slice(f.Hits, func(i, j int) bool {
		return f.Hits[i].Start < f.Hits[j].Start
	})
	return f
}

func (f *FootprintExecution) AggregateTimestamps() *FootprintExecution {
	hitsMap := make(map[int][]int)
	for hitIdx, hit := range f.Hits {
		for _, methodIdx := range hit.Methods {
			hitsList, _ := hitsMap[methodIdx]
			hitsMap[methodIdx] = append(hitsList, hitIdx)
		}
	}

	var uniqueHits []*FootprintExecutionHit
	for _, entries := range hitsMap {
		if len(entries) == 1 {
			uniqueHits = append(uniqueHits, f.Hits[entries[0]])
		} else {
			var start int64
			var end int64
			for _, idx := range entries {
				if start == 0 {
					start = f.Hits[idx].Start
					end = f.Hits[idx].End
				} else {
					if f.Hits[idx].Start < start {
						start = f.Hits[idx].Start
					}
					if f.Hits[idx].End > end {
						end = f.Hits[idx].End
					}
				}
			}
			uniqueHits = append(uniqueHits,
				NewFootprintExecutionHit().
					SetStart(start).
					SetEnd(end).
					SetMethods(f.Hits[entries[0]].Methods[0]).
					SetTestName(f.Hits[entries[0]].TestName).
					SetIsInitFootprints(f.Hits[entries[0]].IsInitFootprints))
		}
	}
	f.Hits = uniqueHits
	return f
}

func (f *FootprintExecution) GetTimestampSpan() (int64, int64) {
	var start int64
	var end int64
	for _, hit := range f.Hits {
		if start == 0 {
			start = hit.Start
			end = hit.End
		} else {
			if hit.Start < start {
				start = hit.Start
			}
			if hit.End > end {
				end = hit.End
			}
		}
	}
	return start, end
}

var _ models.Model = (*FootprintExecution)(nil)
