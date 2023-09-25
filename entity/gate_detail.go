package entity

import "database/sql"

type GateDetail struct {
	GateId      string
	GateName    string
	IsGateUmum  bool
	ExpiredDate sql.NullString
}
