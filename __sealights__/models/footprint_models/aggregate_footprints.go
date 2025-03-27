package footprint_models

import "fmt"

type FootprintAggregationMode int

const (
	FootprintAggregationModeAllHits FootprintAggregationMode = iota
	FootprintAggregationModeCoverage
	FootprintAggregationModeAggregateTimeSpan
)

// Todo: refactor to support aggregation with more than 1 execution time

func AggregateFootprints(footprints []*Footprint, aggType FootprintAggregationMode) (*Footprint, error) {
	if len(footprints) == 0 {
		return nil, fmt.Errorf("no footprints to aggregate")
	}
	aggregatedFootprint := footprints[0].
		updateMethodMap().
		SetAllowEmptyExecutionId(true)
	if err := aggregatedFootprint.Validate(); err != nil {
		return nil, err
	}

	for idx, fp := range footprints[1:] {
		fp.updateMethodMap().
			SetAllowEmptyExecutionId(true)
		if err := fp.Validate(); err != nil {
			return nil, fmt.Errorf("footprint %d validation error: %s", idx+1, err.Error())
		}
		for _, method := range fp.Methods {
			aggregatedFootprint.AddMethod(method)
		}

		for _, hit := range fp.Executions[0].Hits {
			var newMethodsList []int
			for _, methodId := range hit.Methods {
				methodName := fp.Methods[methodId]
				newIdx, ok := aggregatedFootprint.modelMap[methodName]
				if !ok {
					return nil, fmt.Errorf("methodId: %d method: %s - not found", methodId, methodName)
				}
				newMethodsList = append(newMethodsList, newIdx)
			}
			hit.Methods = newMethodsList
			aggregatedFootprint.Executions[0].AddHit(hit)
		}
	}
	return aggregatedFootprint.Complete(aggType), nil
}
