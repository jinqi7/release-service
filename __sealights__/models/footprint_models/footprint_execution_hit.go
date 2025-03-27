package footprint_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type FootprintExecutionHit struct {
	TestName         *string `json:"testName"`
	Start            int64   `json:"start"`
	End              int64   `json:"end"`
	Methods          []int   `json:"methods"`
	IsInitFootprints bool    `json:"isInitFootprints"`
	methodsMap       map[int]bool
}

func NewFootprintExecutionHit() *FootprintExecutionHit {
	return &FootprintExecutionHit{
		methodsMap: make(map[int]bool),
	}
}

func (f *FootprintExecutionHit) Validate() error {
	if f.Start == 0 {
		return fmt.Errorf("FootprintExecutionHit: start is empty")
	}
	if f.End == 0 {
		return fmt.Errorf("FootprintExecutionHit: end is empty")
	}
	if f.Start > f.End {
		return fmt.Errorf("FootprintExecutionHit: start is bigger than end")
	}
	if len(f.Methods) == 0 {
		return fmt.Errorf("FootprintExecutionHit: methods is empty")
	}
	return nil
}

func (f *FootprintExecutionHit) SetTestName(testName *string) *FootprintExecutionHit {
	f.TestName = testName
	return f
}

func (f *FootprintExecutionHit) SetStart(start int64) *FootprintExecutionHit {
	f.Start = start
	return f
}

func (f *FootprintExecutionHit) SetEnd(end int64) *FootprintExecutionHit {
	f.End = end
	return f
}

func (f *FootprintExecutionHit) SetMethods(methods ...int) *FootprintExecutionHit {
	for _, method := range methods {
		if !f.methodsMap[method] {
			f.Methods = append(f.Methods, method)
			f.methodsMap[method] = true
		}
	}
	return f
}

func (f *FootprintExecutionHit) SetIsInitFootprints(isInitFootprints bool) *FootprintExecutionHit {
	f.IsInitFootprints = isInitFootprints
	return f
}

func (f *FootprintExecutionHit) ReplaceMethodIndex(old, new int) *FootprintExecutionHit {
	newMethods := make([]int, 0)
	for _, v := range f.Methods {
		if v == old {
			newMethods = append(newMethods, new)
		} else {
			newMethods = append(newMethods, v)
		}
	}
	f.Methods = newMethods
	return f
}

func (f *FootprintExecutionHit) Key() string {
	if f.TestName == nil {
		return fmt.Sprintf("%d%d", f.Start, f.End)
	} else {
		return fmt.Sprintf("%s%d%d", *f.TestName, f.Start, f.End)
	}
}

func (f *FootprintExecutionHit) Clone() *FootprintExecutionHit {
	newMethodsMap := make(map[int]bool)
	for k, v := range f.methodsMap {
		newMethodsMap[k] = v
	}
	return &FootprintExecutionHit{
		TestName:         f.TestName,
		Start:            f.Start,
		End:              f.End,
		Methods:          f.Methods,
		IsInitFootprints: f.IsInitFootprints,
		methodsMap:       newMethodsMap,
	}
}

var _ models.Model = (*FootprintExecutionHit)(nil)
