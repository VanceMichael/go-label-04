package memory

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"sync"
	"time"
)

type Store struct {
	mu        sync.RWMutex
	users     map[string]domain.User
	byEmail   map[string]string
	sessions  map[string]domain.Session
	shipments map[string]domain.Shipment
	idem      map[string]string
	legs      map[string]domain.FlightLeg
	customs   map[string]domain.CustomsCase
	security  map[string]domain.SecurityCheck
	audit     []domain.AuditEvent
	outbox    map[string]domain.OutboxEvent
}

func New() *Store {
	return &Store{users: map[string]domain.User{}, byEmail: map[string]string{}, sessions: map[string]domain.Session{}, shipments: map[string]domain.Shipment{}, idem: map[string]string{}, legs: map[string]domain.FlightLeg{}, customs: map[string]domain.CustomsCase{}, security: map[string]domain.SecurityCheck{}, outbox: map[string]domain.OutboxEvent{}}
}
func (s *Store) Ping(context.Context) error { return nil }
func (s *Store) GetUserByEmail(_ context.Context, email string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.byEmail[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return cloneUser(s.users[id]), nil
}
func (s *Store) GetUser(_ context.Context, id string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return cloneUser(u), nil
}
func (s *Store) CreateUser(_ context.Context, u domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.byEmail[u.Email]; exists {
		return domain.ErrConflict
	}
	s.users[u.ID] = cloneUser(u)
	s.byEmail[u.Email] = u.ID
	return nil
}
func (s *Store) DeactivateUser(_ context.Context, id string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[id]
	if !ok {
		return domain.ErrNotFound
	}
	u.Active = false
	u.DeactivatedAt = &at
	s.users[id] = u
	for key, session := range s.sessions {
		if session.UserID == id {
			session.RevokedAt = &at
			s.sessions[key] = session
		}
	}
	return nil
}
func (s *Store) CreateSession(_ context.Context, v domain.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[v.TokenHash] = v
	return nil
}
func (s *Store) GetSession(_ context.Context, token string) (domain.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.sessions[tokenHash(token)]
	if !ok {
		return domain.Session{}, domain.ErrNotFound
	}
	return v, nil
}
func (s *Store) RevokeSession(_ context.Context, token string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := tokenHash(token)
	v, ok := s.sessions[key]
	if !ok {
		return domain.ErrNotFound
	}
	v.RevokedAt = &at
	s.sessions[key] = v
	return nil
}
func (s *Store) RevokeUserSessions(_ context.Context, user string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, v := range s.sessions {
		if v.UserID == user {
			v.RevokedAt = &at
			s.sessions[key] = v
		}
	}
	return nil
}
func (s *Store) CreateShipment(_ context.Context, v domain.Shipment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := v.TenantID + "|" + v.IdempotencyKey
	if v.IdempotencyKey != "" {
		if old, ok := s.idem[key]; ok {
			if old == v.ID {
				return nil
			}
			return domain.ErrConflict
		}
	}
	if _, ok := s.shipments[v.ID]; ok {
		return domain.ErrConflict
	}
	s.shipments[v.ID] = v
	if v.IdempotencyKey != "" {
		s.idem[key] = v.ID
	}
	return nil
}
func (s *Store) GetShipment(_ context.Context, tenant, id string) (domain.Shipment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.shipments[id]
	if !ok || v.TenantID != tenant {
		return domain.Shipment{}, domain.ErrNotFound
	}
	return v, nil
}
func (s *Store) UpdateShipment(_ context.Context, v domain.Shipment, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.shipments[v.ID]
	if !ok {
		return domain.ErrNotFound
	}
	if old.Version != version {
		return domain.ErrConflict
	}
	v.Version = version + 1
	s.shipments[v.ID] = v
	return nil
}
func (s *Store) ListShipments(_ context.Context, tenant string, req domain.PageRequest) (domain.Page[domain.Shipment], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	req = req.Normalized()
	all := make([]domain.Shipment, 0)
	for _, v := range s.shipments {
		if v.TenantID == tenant {
			all = append(all, v)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.Before(all[j].CreatedAt) })
	total := len(all)
	if len(all) > req.Limit {
		all = all[:req.Limit]
	}
	return domain.Page[domain.Shipment]{Items: all, Total: total}, nil
}
func (s *Store) FindByIdempotency(_ context.Context, tenant, key string) (domain.Shipment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.idem[tenant+"|"+key]
	if !ok {
		return domain.Shipment{}, domain.ErrNotFound
	}
	return s.shipments[id], nil
}
func (s *Store) CreateLeg(_ context.Context, v domain.FlightLeg) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.legs[v.ID]; ok {
		return domain.ErrConflict
	}
	s.legs[v.ID] = v
	return nil
}
func (s *Store) GetLeg(_ context.Context, tenant, id string) (domain.FlightLeg, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.legs[id]
	if !ok || v.TenantID != tenant {
		return domain.FlightLeg{}, domain.ErrNotFound
	}
	return v, nil
}
func (s *Store) ReserveCapacity(_ context.Context, tenant, id string, weight, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.legs[id]
	if !ok || v.TenantID != tenant {
		return domain.ErrNotFound
	}
	if v.Version != version {
		return domain.ErrConflict
	}
	if err := domain.ValidateCapacity(v, weight); err != nil {
		return err
	}
	v.ReservedKg += weight
	v.Version++
	s.legs[id] = v
	return nil
}
func (s *Store) UpdateLegStatus(_ context.Context, tenant, id string, status domain.LegStatus, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.legs[id]
	if !ok || v.TenantID != tenant {
		return domain.ErrNotFound
	}
	if v.Version != version {
		return domain.ErrConflict
	}
	if !v.Status.CanTransition(status) {
		return domain.ErrState
	}
	v.Status = status
	v.Version++
	s.legs[id] = v
	return nil
}
func (s *Store) BookShipment(_ context.Context, tenant, shipmentID, legID string, shipmentVersion, legVersion int64, at time.Time) (domain.Shipment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	shipment, ok := s.shipments[shipmentID]
	if !ok || shipment.TenantID != tenant {
		return domain.Shipment{}, domain.ErrNotFound
	}
	leg, ok := s.legs[legID]
	if !ok || leg.TenantID != tenant {
		return domain.Shipment{}, domain.ErrNotFound
	}
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
	leg.ReservedKg += shipment.WeightKg
	leg.Version++
	shipment.Status = domain.ShipmentBooked
	shipment.LegID = &legID
	shipment.UpdatedAt = at
	shipment.Version++
	s.legs[legID] = leg
	s.shipments[shipmentID] = shipment
	return shipment, nil
}
func (s *Store) GetCustoms(_ context.Context, id string) (domain.CustomsCase, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.customs[id]
	if !ok {
		return domain.CustomsCase{}, domain.ErrNotFound
	}
	return v, nil
}
func (s *Store) PutCustoms(_ context.Context, v domain.CustomsCase) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.customs[v.ID] = v
	return nil
}
func (s *Store) GetSecurity(_ context.Context, id string) (domain.SecurityCheck, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.security[id]
	if !ok {
		return domain.SecurityCheck{}, domain.ErrNotFound
	}
	return v, nil
}
func (s *Store) PutSecurity(_ context.Context, v domain.SecurityCheck) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.security[v.ID] = v
	return nil
}
func (s *Store) AppendAudit(_ context.Context, v domain.AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audit = append(s.audit, v)
	return nil
}
func (s *Store) ListAudit(_ context.Context, tenant string, req domain.PageRequest) (domain.Page[domain.AuditEvent], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.AuditEvent, 0)
	for _, v := range s.audit {
		if v.TenantID == tenant {
			out = append(out, v)
		}
	}
	return domain.Page[domain.AuditEvent]{Items: out, Total: len(out)}, nil
}
func (s *Store) Enqueue(_ context.Context, v domain.OutboxEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.outbox[v.ID] = v
	return nil
}
func (s *Store) Claim(_ context.Context, now time.Time, n int) ([]domain.OutboxEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.OutboxEvent, 0)
	for id, v := range s.outbox {
		if v.PublishedAt == nil && v.ClaimedAt == nil && !v.AvailableAt.After(now) && len(out) < n {
			v.ClaimedAt = &now
			v.Attempts++
			s.outbox[id] = v
			out = append(out, v)
		}
	}
	return out, nil
}
func (s *Store) MarkPublished(_ context.Context, id string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.outbox[id]
	if !ok {
		return domain.ErrNotFound
	}
	v.PublishedAt = &at
	s.outbox[id] = v
	return nil
}
func (s *Store) MarkFailed(_ context.Context, id string, at time.Time, err string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.outbox[id]
	if !ok {
		return domain.ErrNotFound
	}
	v.LastError = err
	v.ClaimedAt = nil
	v.AvailableAt = at.Add(time.Minute * time.Duration(v.Attempts))
	s.outbox[id] = v
	return nil
}
func (s *Store) Summary(_ context.Context, tenant string) (int, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	pending, failed := 0, 0
	for _, v := range s.outbox {
		if v.TenantID == tenant && v.PublishedAt == nil {
			pending++
			if v.Attempts >= 5 {
				failed++
			}
		}
	}
	return pending, failed, nil
}
func tokenHash(v string) string { sum := sha256.Sum256([]byte(v)); return hex.EncodeToString(sum[:]) }
func cloneUser(v domain.User) domain.User {
	v.PasswordHash = append([]byte(nil), v.PasswordHash...)
	return v
}
