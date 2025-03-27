package buildmap_models

import (
	"fmt"
	"github.com/konflux-ci/release-service/__sealights__/models"
)

type Contributor struct {
	Index int    `json:"-"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewContributor() *Contributor {
	return &Contributor{}
}

func (c *Contributor) SetIndex(index int) *Contributor {
	c.Index = index
	return c
}
func (c *Contributor) SetName(name string) *Contributor {
	c.Name = name
	return c
}

func (c *Contributor) SetEmail(email string) *Contributor {
	c.Email = email
	return c
}

func (c *Contributor) Validate() error {
	if c.Index < 0 {
		return fmt.Errorf("Contributor: index is required")
	}
	if c.Name == "" {
		return fmt.Errorf("Contributor: name is required")
	}
	if c.Email == "" {
		return fmt.Errorf("Contributor: email is required")
	}
	return nil
}

var _ models.Model = (*Contributor)(nil)
