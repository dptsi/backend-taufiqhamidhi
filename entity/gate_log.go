package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	AccessTypeIn  string = "I"
	AccessTypeOut string = "O"
)

type GateLog struct {
	LogId      uuid.UUID
	CardId     string
	GateId     string
	AccessDate time.Time
	AccessType string
}
