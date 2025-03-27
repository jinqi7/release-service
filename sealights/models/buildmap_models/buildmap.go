package buildmap_models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type BuildMap struct {
	Meta         *BuildMapMetadata `json:"meta"`
	Files        []*BuildMapFile   `json:"files"`
	Dependencies []string          `json:"dependencies"`
}

func NewBuildMap() *BuildMap {
	return &BuildMap{
		Dependencies: []string{},
	}
}

func (b *BuildMap) Validate() error {
	if err := b.Meta.Validate(); err != nil {
		return err
	}
	for _, f := range b.Files {
		if err := f.Validate(); err != nil {
			return err
		}
	}
	dupMap := make(map[string]bool)
	for _, f := range b.Files {
		for _, el := range f.stringElements() {
			if _, ok := dupMap[el]; ok {
				return fmt.Errorf("BuildMap: duplicate string element: %s", el)
			} else {
				dupMap[el] = true
			}
		}
	}
	return nil
}

func (b *BuildMap) SetMetadata(meta *BuildMapMetadata) *BuildMap {
	b.Meta = meta
	return b
}

func (b *BuildMap) AddFiles(files ...*BuildMapFile) *BuildMap {
	b.Files = append(b.Files, files...)
	return b
}

func (b *BuildMap) Data() []byte {
	data, _ := json.MarshalIndent(b, "", "  ")
	return data
}

func (b *BuildMap) Load(path string) (*BuildMap, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *BuildMap) AddModuleName(moduleName string) *BuildMap {
	if moduleName == "" {
		return b
	}
	b.Meta.SetUniqueModuleId(moduleName)
	for _, file := range b.Files {
		file.SetLogicalPath(moduleName + "/" + file.LogicalPath)
		for _, method := range file.Methods {
			method.SetUniqueId(fmt.Sprintf("%s.%s", moduleName, method.UniqueId))
			method.SetDisplayName(fmt.Sprintf("%s.%s", moduleName, method.DisplayName))
		}
	}
	return b
}

func (b *BuildMap) MethodsList() []string {
	var methods []string
	for _, f := range b.Files {
		methods = append(methods, f.MethodsList()...)
	}
	return methods
}

var _ models.Model = (*BuildMap)(nil)
