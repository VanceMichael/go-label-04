package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/auth"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/VanceMichael/go-base-airbridge/internal/repository"
	"log/slog"
	"strings"
	"time"
)

type Application struct {
	store repository.Store
	ttl   time.Duration
	log   *slog.Logger
	clock domain.Clock
}

func NewApplication(store repository.Store, ttl time.Duration, logger *slog.Logger) *Application {
	return &Application{store: store, ttl: ttl, log: logger, clock: domain.RealClock{}}
}
func (a *Application) WithClock(clock domain.Clock) *Application { a.clock = clock; return a }
func (a *Application) Ping(ctx context.Context) error            { return a.store.Ping(ctx) }

type LoginResult struct {
	Token     string
	ExpiresAt time.Time
	User      domain.User
}

func (a *Application) Login(ctx context.Context, email, password string) (LoginResult, error) {
	u, err := a.store.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return LoginResult{}, fmt.Errorf("lookup user: %w", err)
	}
	if !u.Active {
		return LoginResult{}, domain.ErrForbidden
	}
	if err := auth.Compare(u.PasswordHash, password); err != nil {
		return LoginResult{}, domain.ErrForbidden
	}
	raw, err := randomToken()
	if err != nil {
		return LoginResult{}, fmt.Errorf("token: %w", err)
	}
	now := a.clock.Now()
	session := domain.Session{ID: randomID(), UserID: u.ID, TokenHash: auth.HashToken(raw), ExpiresAt: now.Add(a.ttl), CreatedAt: now}
	if err := a.store.CreateSession(ctx, session); err != nil {
		return LoginResult{}, fmt.Errorf("create session: %w", err)
	}
	return LoginResult{Token: raw, ExpiresAt: session.ExpiresAt, User: u}, nil
}
func (a *Application) Logout(ctx context.Context, token string) error {
	if token == "" {
		return domain.ErrInvalid
	}
	if err := a.store.RevokeSession(ctx, token, a.clock.Now()); err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}
func (a *Application) Authenticate(ctx context.Context, token string) (domain.User, domain.Session, error) {
	if token == "" {
		return domain.User{}, domain.Session{}, domain.ErrForbidden
	}
	s, err := a.store.GetSession(ctx, token)
	if err != nil {
		return domain.User{}, domain.Session{}, domain.ErrForbidden
	}
	if err := domain.IsSessionActive(s, a.clock.Now()); err != nil {
		return domain.User{}, domain.Session{}, domain.ErrForbidden
	}
	u, err := a.store.GetUser(ctx, s.UserID)
	if err != nil || !u.Active {
		return domain.User{}, domain.Session{}, domain.ErrForbidden
	}
	return u, s, nil
}
func (a *Application) DeactivateUser(ctx context.Context, actor, target string) error {
	u, _, err := a.Authenticate(ctx, actor)
	if err != nil {
		return err
	}
	if u.Role != domain.RoleCoordinator {
		return domain.ErrForbidden
	}
	if err := a.store.DeactivateUser(ctx, target, a.clock.Now()); err != nil {
		return fmt.Errorf("deactivate: %w", err)
	}
	return nil
}

func (a *Application) CreateShipment(ctx context.Context, u domain.User, s domain.Shipment, idem string) (domain.Shipment, error) {
	if u.Role != domain.RoleShipper && u.Role != domain.RoleCoordinator {
		return domain.Shipment{}, domain.ErrForbidden
	}
	s.ID = nonEmpty(s.ID, randomID())
	s.TenantID = u.TenantID
	s.Reference = strings.TrimSpace(s.Reference)
	s.Status = domain.ShipmentDraft
	s.IdempotencyKey = idem
	s.CreatedAt = a.clock.Now()
	s.UpdatedAt = s.CreatedAt
	s.Version = 1
	if err := domain.ValidateShipment(s); err != nil {
		return domain.Shipment{}, err
	}
	if idem != "" {
		old, err := a.store.FindByIdempotency(ctx, u.TenantID, idem)
		if err == nil {
			return old, nil
		}
		if !errors.Is(err, domain.ErrNotFound) {
			return domain.Shipment{}, err
		}
	}
	if err := a.store.CreateShipment(ctx, s); err != nil {
		return domain.Shipment{}, fmt.Errorf("create shipment: %w", err)
	}
	return s, nil
}
func (a *Application) GetShipment(ctx context.Context, u domain.User, id string) (domain.Shipment, error) {
	return a.store.GetShipment(ctx, u.TenantID, id)
}
func (a *Application) ListShipments(ctx context.Context, u domain.User, req domain.PageRequest) (domain.Page[domain.Shipment], error) {
	if u.Role != domain.RoleShipper && u.Role != domain.RoleCoordinator && u.Role != domain.RoleGroundAgent {
		return domain.Page[domain.Shipment]{}, domain.ErrForbidden
	}
	return a.store.ListShipments(ctx, u.TenantID, req)
}
func (a *Application) BookShipment(ctx context.Context, u domain.User, id, legID string) (domain.Shipment, error) {
	if u.Role != domain.RoleCoordinator {
		return domain.Shipment{}, domain.ErrForbidden
	}
	s, err := a.store.GetShipment(ctx, u.TenantID, id)
	if err != nil {
		return domain.Shipment{}, err
	}
	leg, err := a.store.GetLeg(ctx, u.TenantID, legID)
	if err != nil {
		return domain.Shipment{}, err
	}
	if !s.CanTransition(domain.ShipmentBooked) || leg.Status != domain.LegOpen {
		return domain.Shipment{}, domain.ErrState
	}
	if s.Origin != leg.Origin || s.Destination != leg.Destination {
		return domain.Shipment{}, domain.ErrInvalid
	}
	booked, err := a.store.BookShipment(ctx, u.TenantID, id, legID, s.Version, leg.Version, a.clock.Now())
	if err != nil {
		return domain.Shipment{}, fmt.Errorf("book shipment: %w", err)
	}
	return booked, nil
}
func (a *Application) TransitionShipment(ctx context.Context, u domain.User, id string, next domain.ShipmentStatus) (domain.Shipment, error) {
	if u.Role != domain.RoleCoordinator && u.Role != domain.RoleGroundAgent {
		return domain.Shipment{}, domain.ErrForbidden
	}
	s, err := a.store.GetShipment(ctx, u.TenantID, id)
	if err != nil {
		return domain.Shipment{}, err
	}
	if !s.CanTransition(next) {
		return domain.Shipment{}, domain.ErrState
	}
	if s.LegID != nil && (next == domain.ShipmentLoaded || next == domain.ShipmentDeparted) {
		c, ce := a.store.GetCustoms(ctx, id)
		sec, se := a.store.GetSecurity(ctx, id)
		if ce != nil || se != nil || !domain.CustomsAllowsLoading(c) || !domain.SecurityAllowsLoading(sec) {
			return domain.Shipment{}, domain.ErrState
		}
	}
	s.Status = next
	s.UpdatedAt = a.clock.Now()
	if err := a.store.UpdateShipment(ctx, s, s.Version); err != nil {
		return domain.Shipment{}, err
	}
	return s, nil
}

func (a *Application) CreateLeg(ctx context.Context, u domain.User, l domain.FlightLeg) (domain.FlightLeg, error) {
	if u.Role != domain.RoleCoordinator {
		return domain.FlightLeg{}, domain.ErrForbidden
	}
	l.ID = nonEmpty(l.ID, randomID())
	l.TenantID = u.TenantID
	l.Status = domain.LegPlanned
	l.Version = 1
	l.CreatedAt = a.clock.Now()
	if err := domain.ValidateLeg(l, a.clock.Now()); err != nil {
		return domain.FlightLeg{}, err
	}
	if err := a.store.CreateLeg(ctx, l); err != nil {
		return domain.FlightLeg{}, err
	}
	return l, nil
}
func (a *Application) OpenLeg(ctx context.Context, u domain.User, id string) (domain.FlightLeg, error) {
	if u.Role != domain.RoleCoordinator {
		return domain.FlightLeg{}, domain.ErrForbidden
	}
	l, err := a.store.GetLeg(ctx, u.TenantID, id)
	if err != nil {
		return domain.FlightLeg{}, err
	}
	if err := a.store.UpdateLegStatus(ctx, u.TenantID, id, domain.LegOpen, l.Version); err != nil {
		return domain.FlightLeg{}, err
	}
	l.Status = domain.LegOpen
	l.Version++
	return l, nil
}

func (a *Application) PutCustoms(ctx context.Context, u domain.User, c domain.CustomsCase) (domain.CustomsCase, error) {
	if u.Role != domain.RoleCoordinator {
		return domain.CustomsCase{}, domain.ErrForbidden
	}
	if _, err := a.store.GetShipment(ctx, u.TenantID, c.ShipmentID); err != nil {
		return domain.CustomsCase{}, err
	}
	old, err := a.store.GetCustoms(ctx, c.ShipmentID)
	if err == nil && !old.Status.CanTransition(c.Status) {
		return domain.CustomsCase{}, domain.ErrState
	}
	c.ID = nonEmpty(c.ID, randomID())
	c.UpdatedAt = a.clock.Now()
	if err := a.store.PutCustoms(ctx, c); err != nil {
		return domain.CustomsCase{}, err
	}
	return c, nil
}
func (a *Application) PutSecurity(ctx context.Context, u domain.User, s domain.SecurityCheck) (domain.SecurityCheck, error) {
	if u.Role != domain.RoleGroundAgent && u.Role != domain.RoleCoordinator {
		return domain.SecurityCheck{}, domain.ErrForbidden
	}
	if _, err := a.store.GetShipment(ctx, u.TenantID, s.ShipmentID); err != nil {
		return domain.SecurityCheck{}, err
	}
	s.ID = nonEmpty(s.ID, randomID())
	if s.Status == domain.SecurityPassed {
		now := a.clock.Now()
		s.CheckedAt = &now
		s.OfficerID = &u.ID
	}
	if err := a.store.PutSecurity(ctx, s); err != nil {
		return domain.SecurityCheck{}, err
	}
	return s, nil
}
func (a *Application) Summary(ctx context.Context, u domain.User) (domain.OperationsSummary, error) {
	if u.Role != domain.RoleCoordinator {
		return domain.OperationsSummary{}, domain.ErrForbidden
	}
	items, err := a.store.ListShipments(ctx, u.TenantID, domain.PageRequest{Limit: 200})
	if err != nil {
		return domain.OperationsSummary{}, err
	}
	summary := domain.OperationsSummary{TenantID: u.TenantID}
	for _, s := range items.Items {
		switch s.Status {
		case domain.ShipmentDraft:
			summary.Draft++
		case domain.ShipmentBooked, domain.ShipmentScreening, domain.ShipmentCleared, domain.ShipmentLoaded:
			summary.Booked++
		case domain.ShipmentDeparted:
			summary.InFlight++
		}
	}
	pending, failed, err := a.store.Summary(ctx, u.TenantID)
	if err != nil {
		return domain.OperationsSummary{}, err
	}
	summary.PendingOutbox = pending
	summary.FailedOutbox = failed
	return summary, nil
}

func randomID() string { b := make([]byte, 16); _, _ = rand.Read(b); return hex.EncodeToString(b) }
func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func nonEmpty(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
