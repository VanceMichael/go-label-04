package postgres

import (
	"context"
	"errors"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (s *Store) GetCustoms(ctx context.Context, id string) (domain.CustomsCase, error) {
	var v domain.CustomsCase
	var st string
	err := s.db.QueryRow(ctx, `SELECT id,shipment_id,status,document_ref,reviewed_by,updated_at FROM customs_cases WHERE shipment_id=$1`, id).Scan(&v.ID, &v.ShipmentID, &st, &v.DocumentRef, &v.ReviewedBy, &v.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, domain.ErrNotFound
	}
	v.Status = domain.CustomsStatus(st)
	return v, err
}
func (s *Store) PutCustoms(ctx context.Context, v domain.CustomsCase) error {
	_, err := s.db.Exec(ctx, `INSERT INTO customs_cases(id,shipment_id,status,document_ref,reviewed_by,updated_at) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT(shipment_id) DO UPDATE SET status=EXCLUDED.status,document_ref=EXCLUDED.document_ref,reviewed_by=EXCLUDED.reviewed_by,updated_at=EXCLUDED.updated_at`, v.ID, v.ShipmentID, v.Status, v.DocumentRef, v.ReviewedBy, v.UpdatedAt)
	return err
}
func (s *Store) GetSecurity(ctx context.Context, id string) (domain.SecurityCheck, error) {
	var v domain.SecurityCheck
	var st string
	err := s.db.QueryRow(ctx, `SELECT id,shipment_id,status,officer_id,notes,checked_at FROM security_checks WHERE shipment_id=$1`, id).Scan(&v.ID, &v.ShipmentID, &st, &v.OfficerID, &v.Notes, &v.CheckedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, domain.ErrNotFound
	}
	v.Status = domain.SecurityStatus(st)
	return v, err
}
func (s *Store) PutSecurity(ctx context.Context, v domain.SecurityCheck) error {
	_, err := s.db.Exec(ctx, `INSERT INTO security_checks(id,shipment_id,status,officer_id,notes,checked_at) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT(shipment_id) DO UPDATE SET status=EXCLUDED.status,officer_id=EXCLUDED.officer_id,notes=EXCLUDED.notes,checked_at=EXCLUDED.checked_at`, v.ID, v.ShipmentID, v.Status, v.OfficerID, v.Notes, v.CheckedAt)
	return err
}
