package repository

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"time"
)

type UserRepository interface {
	GetUserByEmail(context.Context, string) (domain.User, error)
	GetUser(context.Context, string) (domain.User, error)
	CreateUser(context.Context, domain.User) error
	DeactivateUser(context.Context, string, time.Time) error
}
type SessionRepository interface {
	CreateSession(context.Context, domain.Session) error
	GetSession(context.Context, string) (domain.Session, error)
	RevokeSession(context.Context, string, time.Time) error
	RevokeUserSessions(context.Context, string, time.Time) error
}
type ShipmentRepository interface {
	CreateShipment(context.Context, domain.Shipment) error
	GetShipment(context.Context, string, string) (domain.Shipment, error)
	UpdateShipment(context.Context, domain.Shipment, int64) error
	ListShipments(context.Context, string, domain.PageRequest) (domain.Page[domain.Shipment], error)
	FindByIdempotency(context.Context, string, string) (domain.Shipment, error)
}
type LegRepository interface {
	CreateLeg(context.Context, domain.FlightLeg) error
	GetLeg(context.Context, string, string) (domain.FlightLeg, error)
	ReserveCapacity(context.Context, string, string, int64, int64) error
	UpdateLegStatus(context.Context, string, string, domain.LegStatus, int64) error
}
type BookingRepository interface {
	BookShipment(context.Context, string, string, string, int64, int64, time.Time) (domain.Shipment, error)
}
type ComplianceRepository interface {
	GetCustoms(context.Context, string) (domain.CustomsCase, error)
	PutCustoms(context.Context, domain.CustomsCase) error
	GetSecurity(context.Context, string) (domain.SecurityCheck, error)
	PutSecurity(context.Context, domain.SecurityCheck) error
}
type AuditRepository interface {
	AppendAudit(context.Context, domain.AuditEvent) error
	ListAudit(context.Context, string, domain.PageRequest) (domain.Page[domain.AuditEvent], error)
}
type OutboxRepository interface {
	Enqueue(context.Context, domain.OutboxEvent) error
	Claim(context.Context, time.Time, int) ([]domain.OutboxEvent, error)
	MarkPublished(context.Context, string, time.Time) error
	MarkFailed(context.Context, string, time.Time, string) error
	Summary(context.Context, string) (int, int, error)
}
type Store interface {
	UserRepository
	SessionRepository
	ShipmentRepository
	LegRepository
	BookingRepository
	ComplianceRepository
	AuditRepository
	OutboxRepository
	Ping(context.Context) error
}
