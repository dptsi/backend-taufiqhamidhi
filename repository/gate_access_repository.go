package repository

import "myits-gate-api/entity"

type GateAccessRepository interface {
	FindBy(gateId string, cardId string) (*entity.GateAccess, error)
	FindByUmum(cardId string) (*entity.GateAccess, error)
	GetGate(gateId string) (*entity.GateDetail, error)
}
