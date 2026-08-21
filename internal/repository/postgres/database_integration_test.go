package postgres

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestPostgresMigrationsAndRecovery(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := Open(ctx, EnvURL())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	store := NewStore(db)
	if err := store.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	tenant := "integration-tenant"
	_, _ = db.Pool.Exec(ctx, `DELETE FROM shipments WHERE tenant_id=$1`, tenant)
	_, _ = db.Pool.Exec(ctx, `DELETE FROM users WHERE tenant_id=$1`, tenant)
	_, err = db.Pool.Exec(ctx, `INSERT INTO tenants(id,name,active,created_at) VALUES($1,$2,true,$3) ON CONFLICT(id) DO NOTHING`, tenant, "Integration", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	shipment := domain.Shipment{ID: "integration-shipment", TenantID: tenant, Reference: "INT-1", Origin: "PEK", Destination: "FRA", WeightKg: 10, Pieces: 1, Status: domain.ShipmentDraft, Version: 1, CreatedAt: now, UpdatedAt: now}
	if err := store.CreateShipment(ctx, shipment); err != nil {
		t.Fatal(err)
	}
	got, err := store.GetShipment(ctx, tenant, shipment.ID)
	if err != nil || got.Reference != "INT-1" {
		t.Fatalf("%v %#v", err, got)
	}
	db.Close()
	db, err = Open(ctx, EnvURL())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	got, err = NewStore(db).GetShipment(ctx, tenant, shipment.ID)
	if err != nil || got.ID != shipment.ID {
		t.Fatalf("recovery %v %#v", err, got)
	}
}
