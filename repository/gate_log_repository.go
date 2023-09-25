package repository

import "myits-gate-api/entity"

type GateLogRepository interface {
	Save(log entity.GateLog) error
}
