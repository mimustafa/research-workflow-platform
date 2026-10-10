package workflow

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrEmptyBrief = errors.New("brief must not be empty")

type Status string

const (
	StatusQueued Status = "queued"
)

type Workflow struct {
	ID        string    `json:"id"`
	Brief     string    `json:"brief"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func New(brief string, now time.Time) (Workflow, error) {
	brief = strings.TrimSpace(brief)
	if brief == "" {
		return Workflow{}, ErrEmptyBrief
	}

	return Workflow{
		ID:        uuid.NewString(),
		Brief:     brief,
		Status:    StatusQueued,
		CreatedAt: now.UTC(),
	}, nil
}
