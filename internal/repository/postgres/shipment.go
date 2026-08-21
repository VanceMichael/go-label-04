package postgres

import (
	"context"
	"errors"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/jackc/pgx/v5"
	"time"
)

func (s *Store) CreateShipment(ctx context.Context, v domain.Shipment) error {
	_, err := s.db.Exec(ctx, `INSERT INTO shipments(id,tenant_id,reference,origin,destination,weight_kg,pieces,status,leg_id,idempotency_key,version,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`, v.ID, v.TenantID, v.Reference, v.Origin, v.Destination, v.WeightKg, v.Pieces, v.Status, v.LegID, v.IdempotencyKey, v.Version, v.CreatedAt, v.UpdatedAt)
	if isUnique(err) {
		return domain.ErrConflict
	}
	return err
}
func (s *Store) GetShipment(ctx context.Context, tenant, id string) (domain.Shipment, error) {
	return s.shipment(ctx, `SELECT id,tenant_id,reference,origin,destination,weight_kg,pieces,status,leg_id,idempotency_key,version,created_at,updated_at FROM shipments WHERE tenant_id=$1 AND id=$2`, tenant, id)
}
func (s *Store) shipment(ctx context.Context, q string, args ...any) (domain.Shipment, error) {
	var v domain.Shipment
	var status string
	err := s.db.QueryRow(ctx, q, args...).Scan(&v.ID, &v.TenantID, &v.Reference, &v.Origin, &v.Destination, &v.WeightKg, &v.Pieces, &status, &v.LegID, &v.IdempotencyKey, &v.Version, &v.CreatedAt, &v.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, domain.ErrNotFound
	}
	v.Status = domain.ShipmentStatus(status)
	return v, err
}
func (s *Store) UpdateShipment(ctx context.Context, v domain.Shipment, version int64) error {
	tag, err := s.db.Exec(ctx, `UPDATE shipments SET reference=$3,origin=$4,destination=$5,weight_kg=$6,pieces=$7,status=$8,leg_id=$9,updated_at=$10,version=version+1 WHERE tenant_id=$1 AND id=$2 AND version=$11`, v.TenantID, v.ID, v.Reference, v.Origin, v.Destination, v.WeightKg, v.Pieces, v.Status, v.LegID, v.UpdatedAt, version)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrConflict
	}
	return nil
}
func (s *Store) ListShipments(ctx context.Context, tenant string, req domain.PageRequest) (domain.Page[domain.Shipment], error) {
	req = req.Normalized()
	rows, err := s.db.Query(ctx, `SELECT id,tenant_id,reference,origin,destination,weight_kg,pieces,status,leg_id,idempotency_key,version,created_at,updated_at FROM shipments WHERE tenant_id=$1 ORDER BY created_at,id LIMIT $2`, tenant, req.Limit)
	if err != nil {
		return domain.Page[domain.Shipment]{}, err
	}
	defer rows.Close()
	out := make([]domain.Shipment, 0)
	for rows.Next() {
		var v domain.Shipment
		var st string
		if err := rows.Scan(&v.ID, &v.TenantID, &v.Reference, &v.Origin, &v.Destination, &v.WeightKg, &v.Pieces, &st, &v.LegID, &v.IdempotencyKey, &v.Version, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return domain.Page[domain.Shipment]{}, err
		}
		v.Status = domain.ShipmentStatus(st)
		out = append(out, v)
	}
	var total int
	if err := s.db.QueryRow(ctx, `SELECT count(*) FROM shipments WHERE tenant_id=$1`, tenant).Scan(&total); err != nil {
		return domain.Page[domain.Shipment]{}, err
	}
	return domain.Page[domain.Shipment]{Items: out, Total: total}, rows.Err()
}
func (s *Store) FindByIdempotency(ctx context.Context, tenant, key string) (domain.Shipment, error) {
	return s.shipment(ctx, `SELECT id,tenant_id,reference,origin,destination,weight_kg,pieces,status,leg_id,idempotency_key,version,created_at,updated_at FROM shipments WHERE tenant_id=$1 AND idempotency_key=$2`, tenant, key)
}
func (s *Store) CreateLeg(ctx context.Context, v domain.FlightLeg) error {
	_, err := s.db.Exec(ctx, `INSERT INTO flight_legs(id,tenant_id,flight_number,origin,destination,departure_at,arrival_at,capacity_kg,reserved_kg,status,version,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, v.ID, v.TenantID, v.FlightNumber, v.Origin, v.Destination, v.DepartureAt, v.ArrivalAt, v.CapacityKg, v.ReservedKg, v.Status, v.Version, v.CreatedAt)
	return err
}
func (s *Store) GetLeg(ctx context.Context, tenant, id string) (domain.FlightLeg, error) {
	var v domain.FlightLeg
	var st string
	err := s.db.QueryRow(ctx, `SELECT id,tenant_id,flight_number,origin,destination,departure_at,arrival_at,capacity_kg,reserved_kg,status,version,created_at FROM flight_legs WHERE tenant_id=$1 AND id=$2`, tenant, id).Scan(&v.ID, &v.TenantID, &v.FlightNumber, &v.Origin, &v.Destination, &v.DepartureAt, &v.ArrivalAt, &v.CapacityKg, &v.ReservedKg, &st, &v.Version, &v.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, domain.ErrNotFound
	}
	v.Status = domain.LegStatus(st)
	return v, err
}
func (s *Store) ReserveCapacity(ctx context.Context, tenant, id string, weight, version int64) error {
	tag, err := s.db.Exec(ctx, `UPDATE flight_legs SET reserved_kg=reserved_kg+$4,version=version+1 WHERE tenant_id=$1 AND id=$2 AND version=$3 AND reserved_kg+$4<=capacity_kg`, tenant, id, version, weight)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrConflict
	}
	return nil
}
func (s *Store) UpdateLegStatus(ctx context.Context, tenant, id string, status domain.LegStatus, version int64) error {
	tag, err := s.db.Exec(ctx, `UPDATE flight_legs SET status=$4,version=version+1 WHERE tenant_id=$1 AND id=$2 AND version=$3`, tenant, id, version, status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrConflict
	}
	return nil
}

var _ = time.UTC
