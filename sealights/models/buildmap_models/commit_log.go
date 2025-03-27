package buildmap_models

import (
	"fmt"
	"github.com/konflux-ci/release-service/__sealights__/models"
)

type CommitLog struct {
	Commit           string `json:"commit"`
	AuthorDate       int64  `json:"authorDate"`
	CommitterDate    int64  `json:"committerDate"`
	Title            string `json:"title"`
	ContributorIndex int    `json:"contributorIndex"`
	ContributorName  string `json:"-"`
	ContributorEmail string `json:"-"`
}

func NewCommitLog() *CommitLog {
	return &CommitLog{}
}

func (c *CommitLog) SetContributorName(name string) *CommitLog {
	c.ContributorName = name
	return c
}

func (c *CommitLog) SetCommit(commit string) *CommitLog {
	c.Commit = commit
	return c
}

func (c *CommitLog) SetAuthorDate(authorDate int64) *CommitLog {
	c.AuthorDate = authorDate
	return c
}

func (c *CommitLog) SetCommitterDate(committerDate int64) *CommitLog {
	c.CommitterDate = committerDate
	return c
}

func (c *CommitLog) SetTitle(title string) *CommitLog {
	c.Title = title
	return c
}

func (c *CommitLog) SetContributorIndex(contributorIndex int) *CommitLog {
	c.ContributorIndex = contributorIndex
	return c
}

func (c *CommitLog) SetContributorEmail(email string) *CommitLog {
	c.ContributorEmail = email
	return c
}

func (c *CommitLog) Validate() error {
	if c.Commit == "" {
		return fmt.Errorf("CommitLog: commit is required")
	}
	if c.AuthorDate == 0 {
		return fmt.Errorf("CommitLog: authorDate is required")
	}
	if c.CommitterDate == 0 {
		return fmt.Errorf("CommitLog: committerDate is required")
	}
	if c.Title == "" {
		return fmt.Errorf("CommitLog: title is required")
	}
	if c.ContributorIndex < 0 {
		return fmt.Errorf("CommitLog: contributorIndex is required")
	}
	return nil
}

func (c *CommitLog) String() string {
	return fmt.Sprintf("CommitLog{Commit: %s, AuthorDate: %d, CommitterDate: %d, Title: %s, ContributorIndex: %d, ContributorName: %s, ContributorEmail: %s}",
		c.Commit, c.AuthorDate, c.CommitterDate, c.Title, c.ContributorIndex, c.ContributorName, c.ContributorEmail)
}

var _ models.Model = (*CommitLog)(nil)
