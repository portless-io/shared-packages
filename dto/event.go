package dto

import (
	"time"

	"github.com/portless-io/shared-packages/errors"
)

type CreateEvent struct {
	Type      string      `json:"type"`
	Resource  string      `json:"resource"`
	Data      interface{} `json:"data"`
	CreatedBy string      `json:"createdBy"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

func (createEvent *CreateEvent) Validate() error {
	if createEvent.CreatedBy == "" {
		return errors.NewInvalidArgumentErr("field createdBy required")
	}

	if createEvent.Type == "" {
		return errors.NewInvalidArgumentErr("field type required")
	}

	if createEvent.Resource == "" {
		return errors.NewInvalidArgumentErr("field resource required")
	}

	return nil
}
