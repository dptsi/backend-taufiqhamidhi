package entity

import "time"

type GateAccess struct {
	CardId      string
	GateId      string
	ExpiredDate time.Time
	Active      bool
}
