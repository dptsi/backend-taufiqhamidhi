package repository

import (
	"context"
	"database/sql"
	"log"
	"myits-gate-api/entity"
)

type SqlServerGateLogRepository struct {
	db  *sql.DB
	ctx context.Context
}

func NewSqlServerGateLogRepository(db *sql.DB, ctx context.Context) GateLogRepository {
	return &SqlServerGateLogRepository{db, ctx}
}

func (s *SqlServerGateLogRepository) Save(l entity.GateLog) error {

	tsql := `INSERT INTO log_gate (id_log_gate, id_kartu_akses, id_gate, tgl_akses, jenis_akses, created_at, updated_at, deleted_at, updater)
			 VALUES (@LogId, @CardId, @GateId, @AccessDate, @AccessType, GETDATE(), GETDATE(), NULL, NULL)`

	stmt, err := s.db.Prepare(tsql)

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		s.ctx,
		sql.Named("LogId", l.LogId),
		sql.Named("CardId", l.CardId),
		sql.Named("GateId", l.GateId),
		sql.Named("AccessDate", l.AccessDate),
		sql.Named("AccessType", l.AccessType),
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
