package postgres

import (
	"context"
	"errors"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/jackc/pgx/v5"
	"time"
)

func (s *Store) BookShipment(ctx context.Context, tenant, shipmentID, legID string, shipmentVersion, legVersion int64, at time.Time) (domain.Shipment, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.Shipment{}, err
	}
	defer tx.Rollback(ctx)
	var shipment domain.Shipment
	var shipmentStatus string
	err = tx.QueryRow(ctx, `SELECT id,tenant_id,reference,origin,destination,weight_kg,pieces,status,leg_id,idempotency_key,version,created_at,updated_at FROM shipments WHERE tenant_id=$1 AND id=$2 FOR UPDATE`, tenant, shipmentID).Scan(&shipment.ID, &shipment.TenantID, &shipment.Reference, &shipment.Origin, &shipment.Destination, &shipment.WeightKg, &shipment.Pieces, &shipmentStatus, &shipment.LegID, &shipment.IdempotencyKey, &shipment.Version, &shipment.CreatedAt, &shipment.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Shipment{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Shipment{}, err
	}
	shipment.Status = domain.ShipmentStatus(shipmentStatus)
	var leg domain.FlightLeg
	var legStatus string
	err = tx.QueryRow(ctx, `SELECT id,tenant_id,flight_number,origin,destination,departure_at,arrival_at,capacity_kg,reserved_kg,status,version,created_at FROM flight_legs WHERE tenant_id=$1 AND id=$2 FOR UPDATE`, tenant, legID).Scan(&leg.ID, &leg.TenantID, &leg.FlightNumber, &leg.Origin, &leg.Destination, &leg.DepartureAt, &leg.ArrivalAt, &leg.CapacityKg, &leg.ReservedKg, &legStatus, &leg.Version, &leg.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Shipment{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Shipment{}, err
	}
	leg.Status = domain.LegStatus(legStatus)
	if shipment.Version != shipmentVersion || leg.Version != legVersion {
		return domain.Shipment{}, domain.ErrConflict
	}
	if !shipment.CanTransition(domain.ShipmentBooked) || leg.Status != domain.LegOpen {
		return domain.Shipment{}, domain.ErrState
	}
	if shipment.Origin != leg.Origin || shipment.Destination != leg.Destination {
		return domain.Shipment{}, domain.ErrInvalid
	}
	if err := domain.ValidateCapacity(leg, shipment.WeightKg); err != nil {
		return domain.Shipment{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE flight_legs SET reserved_kg=reserved_kg+$3,version=version+1 WHERE tenant_id=$1 AND id=$2`, tenant, legID, shipment.WeightKg); err != nil {
		return domain.Shipment{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE shipments SET status=$3,leg_id=$4,updated_at=$5,version=version+1 WHERE tenant_id=$1 AND id=$2`, tenant, shipmentID, domain.ShipmentBooked, legID, at); err != nil {
		return domain.Shipment{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Shipment{}, err
	}
	shipment.Status = domain.ShipmentBooked
	shipment.LegID = &legID
	shipment.UpdatedAt = at
	shipment.Version++
	return shipment, nil
}
