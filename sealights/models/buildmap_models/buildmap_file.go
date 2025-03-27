package buildmap_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type BuildMapFile struct {
	LogicalPath  string            `json:"logicalPath"`
	PhysicalPath string            `json:"physicalPath"`
	Methods      []*BuildMapMethod `json:"methods"`
}

func NewBuildMapFile() *BuildMapFile {
	return &BuildMapFile{}
}

func (b *BuildMapFile) SetLogicalPath(logicalPath string) *BuildMapFile {
	b.LogicalPath = logicalPath
	return b
}

func (b *BuildMapFile) SetPhysicalPath(physicalPath string) *BuildMapFile {
	b.PhysicalPath = physicalPath
	return b
}

func (b *BuildMapFile) AddMethod(method *BuildMapMethod) *BuildMapFile {
	b.Methods = append(b.Methods, method)
	return b
}

func (b *BuildMapFile) Validate() error {
	if b.LogicalPath == "" {
		return fmt.Errorf("BuildMapFile: logicalPath is required")
	}
	if b.PhysicalPath == "" {
		return fmt.Errorf("BuildMapFile: physicalPath is required")
	}
	if len(b.Methods) == 0 {
		return fmt.Errorf("BuildMapFile: methods is required")
	}
	for _, method := range b.Methods {
		if err := method.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (b *BuildMapFile) stringElements() []string {
	el := []string{
		b.LogicalPath,
	}
	for _, method := range b.Methods {
		el = append(el, method.stringElements()...)
	}
	return el
}

func (b *BuildMapFile) MethodsList() []string {
	var methods []string
	for _, method := range b.Methods {
		methods = append(methods, method.UniqueId)
	}
	return methods
}

var _ models.Model = (*BuildMapFile)(nil)
