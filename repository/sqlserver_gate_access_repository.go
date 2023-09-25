package repository

import (
	"context"
	"database/sql"
	"log"
	"myits-gate-api/entity"
	"time"
)

type SqlServerGateAccessRepository struct {
	db  *sql.DB
	ctx context.Context
}

func NewSqlServerGateAccessRepository(db *sql.DB, ctx context.Context) GateAccessRepository {
	return &SqlServerGateAccessRepository{db, ctx}
}

func (s *SqlServerGateAccessRepository) FindBy(gateId string, cardId string) (*entity.GateAccess, error) {

	tsql := `SELECT id_kartu_akses, id_gate, tgl_kadaluarsa, is_aktif
			 FROM dbo.akses_gate a
			 WHERE a.id_gate = @GateId AND a.id_kartu_akses = @CardId AND a.deleted_at IS NULL`

	rows, err := s.db.QueryContext(s.ctx, tsql, sql.Named("GateId", gateId), sql.Named("CardId", cardId))
	log.Println(rows)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {

		var (
			cardId      string
			gateId      string
			expiredDate time.Time
			active      bool
		)

		err = rows.Scan(&cardId, &gateId, &expiredDate, &active)

		if err != nil {
			return nil, err
		}

		return &entity.GateAccess{
			CardId:      cardId,
			GateId:      gateId,
			ExpiredDate: expiredDate,
			Active:      active,
		}, nil

	}

	return nil, nil

}

func (s *SqlServerGateAccessRepository) FindByUmum(cardId string) (*entity.GateAccess, error) {

	tsql := `SELECT id_kartu_akses, tgl_kadaluarsa, is_aktif
			 FROM dbo.kartu_akses a
			 WHERE a.id_kartu_akses = @CardId AND a.deleted_at IS NULL`

	rows, err := s.db.QueryContext(s.ctx, tsql, sql.Named("CardId", cardId))
	log.Println(rows)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {

		var (
			cardId      string
			expiredDate time.Time
			active      bool
		)

		err = rows.Scan(&cardId, &expiredDate, &active)

		if err != nil {
			return nil, err
		}

		return &entity.GateAccess{
			CardId:      cardId,
			ExpiredDate: expiredDate,
			Active:      active,
		}, nil

	}

	return nil, nil

}

func (s *SqlServerGateAccessRepository) GetGate(gateId string) (*entity.GateDetail, error) {

	tsql := `SELECT TOP 100 id_gate, nama, is_gate_umum, updated_at FROM ref.gate g WHERE g.id_gate = @GateId AND g.expired_at IS NULL`

	rows, err := s.db.QueryContext(s.ctx, tsql, sql.Named("GateId", gateId))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {

		var (
			gateId      string
			gateName    string
			isGateUmum  bool
			expiredDate sql.NullString
		)

		err = rows.Scan(&gateId, &gateName, &isGateUmum, &expiredDate)

		if err != nil {
			return nil, err
		}
		return &entity.GateDetail{
			GateId:      gateId,
			GateName:    gateName,
			IsGateUmum:  isGateUmum,
			ExpiredDate: expiredDate,
		}, nil

	}

	return nil, nil

}
