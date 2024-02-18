package status

import (
	"github.com/retr0h/psion/pkg/resource/api"
)

// GetType the action property.
func (sc *StatusConditions) GetType() api.SpecAction { return sc.Type }

// GetTypeString the action property.
func (sc *StatusConditions) GetTypeString() string { return string(sc.GetType()) }

// GetStatus the status property.
func (sc *StatusConditions) GetStatus() api.Phase { return sc.Status }

// GetMessage get the message property.
func (sc *StatusConditions) GetMessage() string { return sc.Message }

// SetMessage set the message property.
func (sc *StatusConditions) SetMessage(message string) { sc.Message = message }

// GetReason the reason property.
func (sc *StatusConditions) GetReason() api.Action { return sc.Reason }

// GetReasonString the reason property.
func (sc *StatusConditions) GetReasonString() string { return string(sc.GetReason()) }

// GetStatusString the status property as a string.
func (sc *StatusConditions) GetStatusString() string { return string(sc.GetStatus()) }

// GetGot get the got property.
func (sc *StatusConditions) GetGot() string { return sc.Got }

// GetWant get the want property.
func (sc *StatusConditions) GetWant() string { return sc.Want }
