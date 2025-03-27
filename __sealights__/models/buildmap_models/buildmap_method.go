package buildmap_models

import (
	"fmt"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type BuildMapMethod struct {
	UniqueId  string `json:"uniqueId"`
	Hash      string `json:"hash"`
	ElementId string `json:"elementId"`
	Meta      struct {
		Type      string `json:"type"`
		Anonymous bool   `json:"anonymous"`
	} `json:"meta"`
	DisplayName string `json:"displayName"`
	Position    []int  `json:"position"`
	EndPosition []int  `json:"endPosition"`
}

func NewBuildMapMethod() *BuildMapMethod {
	return &BuildMapMethod{}
}

func (b *BuildMapMethod) SetUniqueId(uniqueId string) *BuildMapMethod {
	b.UniqueId = uniqueId
	return b
}

func (b *BuildMapMethod) SetHash(hash string) *BuildMapMethod {
	b.Hash = hash
	return b
}

func (b *BuildMapMethod) SetElementId(elementId string) *BuildMapMethod {
	b.ElementId = elementId
	return b
}

func (b *BuildMapMethod) SetMeta(metaType string, isAnonymous bool) *BuildMapMethod {
	b.Meta = struct {
		Type      string `json:"type"`
		Anonymous bool   `json:"anonymous"`
	}{
		Type:      metaType,
		Anonymous: isAnonymous,
	}
	return b
}

func (b *BuildMapMethod) SetDisplayName(displayName string) *BuildMapMethod {
	b.DisplayName = displayName
	return b
}

func (b *BuildMapMethod) SetPosition(position []int) *BuildMapMethod {
	b.Position = position
	return b
}

func (b *BuildMapMethod) SetEndPosition(endPosition []int) *BuildMapMethod {
	b.EndPosition = endPosition
	return b
}

func (b *BuildMapMethod) Validate() error {
	if b.UniqueId == "" {
		return fmt.Errorf("BuildMapMethod: uniqueId is required")
	}
	if b.Hash == "" {
		return fmt.Errorf("BuildMapMethod: hash is required")
	}
	if b.ElementId == "" {
		return fmt.Errorf("BuildMapMethod: elementId is required")
	}
	if b.Meta.Type == "" {
		return fmt.Errorf("BuildMapMethod: type is required")
	}
	if b.DisplayName == "" {
		return fmt.Errorf("BuildMapMethod: displayName is required")
	}
	if len(b.Position) != 2 {
		return fmt.Errorf("BuildMapMethod: position must be an array of length 2")
	}
	if len(b.EndPosition) != 2 {
		return fmt.Errorf("BuildMapMethod: endPosition must be an array of length 2")
	}
	return nil
}

func (b *BuildMapMethod) stringElements() []string {
	return []string{
		b.UniqueId,
		b.Hash,
		// b.DisplayName,
	}
}

var _ models.Model = (*BuildMapMethod)(nil)
