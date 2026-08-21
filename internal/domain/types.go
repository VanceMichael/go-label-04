package domain

import "time"

type Role string

const (
	RoleShipper     Role = "shipper"
	RoleCoordinator Role = "coordinator"
	RoleGroundAgent Role = "ground_agent"
)

type ShipmentStatus string

const (
	ShipmentDraft     ShipmentStatus = "draft"
	ShipmentBooked    ShipmentStatus = "booked"
	ShipmentScreening ShipmentStatus = "screening"
	ShipmentCleared   ShipmentStatus = "cleared"
	ShipmentLoaded    ShipmentStatus = "loaded"
	ShipmentDeparted  ShipmentStatus = "departed"
	ShipmentCancelled ShipmentStatus = "cancelled"
)

type LegStatus string

const (
	LegPlanned  LegStatus = "planned"
	LegOpen     LegStatus = "open"
	LegBoarding LegStatus = "boarding"
	LegDeparted LegStatus = "departed"
	LegClosed   LegStatus = "closed"
)

type CustomsStatus string

const (
	CustomsPending  CustomsStatus = "pending"
	CustomsReview   CustomsStatus = "review"
	CustomsReleased CustomsStatus = "released"
	CustomsHeld     CustomsStatus = "held"
)

type SecurityStatus string

const (
	SecurityPending SecurityStatus = "pending"
	SecurityPassed  SecurityStatus = "passed"
	SecurityFailed  SecurityStatus = "failed"
)

type Tenant struct {
	ID        string
	Name      string
	Active    bool
	CreatedAt time.Time
}
type User struct {
	ID            string
	TenantID      string
	Email         string
	PasswordHash  []byte
	Role          Role
	Active        bool
	CreatedAt     time.Time
	DeactivatedAt *time.Time
}
type Session struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
type Shipment struct {
	ID             string
	TenantID       string
	Reference      string
	Origin         string
	Destination    string
	WeightKg       int64
	Pieces         int
	Status         ShipmentStatus
	LegID          *string
	IdempotencyKey string
	Version        int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
type FlightLeg struct {
	ID           string
	TenantID     string
	FlightNumber string
	Origin       string
	Destination  string
	DepartureAt  time.Time
	ArrivalAt    time.Time
	CapacityKg   int64
	ReservedKg   int64
	Status       LegStatus
	Version      int64
	CreatedAt    time.Time
}
type CustomsCase struct {
	ID          string
	ShipmentID  string
	Status      CustomsStatus
	DocumentRef string
	ReviewedBy  *string
	UpdatedAt   time.Time
}
type SecurityCheck struct {
	ID         string
	ShipmentID string
	Status     SecurityStatus
	OfficerID  *string
	Notes      string
	CheckedAt  *time.Time
}
type HandlingEvent struct {
	ID         string
	ShipmentID string
	Kind       string
	ActorID    string
	OccurredAt time.Time
	Metadata   map[string]string
}
type AuditEvent struct {
	ID         string
	TenantID   string
	ActorID    string
	ObjectType string
	ObjectID   string
	Action     string
	Result     string
	RequestID  string
	OccurredAt time.Time
}
type OutboxEvent struct {
	ID          string
	TenantID    string
	Topic       string
	AggregateID string
	Payload     []byte
	Attempts    int
	AvailableAt time.Time
	ClaimedAt   *time.Time
	PublishedAt *time.Time
	LastError   string
}
type OperationsSummary struct {
	TenantID      string
	Draft         int
	Booked        int
	InFlight      int
	Held          int
	PendingOutbox int
	FailedOutbox  int
}
