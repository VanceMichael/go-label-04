package scan

import (
	"context"
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
	"time"
)

type Result struct {
	ShipmentID string
	Status     domain.SecurityStatus
	Officer    string
	Reason     string
	CheckedAt  time.Time
}
type Service struct {
	mu      sync.RWMutex
	results map[string]Result
	blocked map[string]string
}

func New() *Service { return &Service{results: map[string]Result{}, blocked: map[string]string{}} }
func (s *Service) Check(ctx context.Context, shipment, officer string, signals map[string]bool) (Result, error) {
	select {
	case <-ctx.Done():
		return Result{}, ctx.Err()
	default:
	}
	if shipment == "" || officer == "" {
		return Result{}, domain.ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for name, bad := range signals {
		if bad {
			r := Result{ShipmentID: shipment, Officer: officer, Status: domain.SecurityFailed, Reason: fmt.Sprintf("signal %s", name), CheckedAt: time.Now().UTC()}
			s.results[shipment] = r
			s.blocked[shipment] = r.Reason
			return r, nil
		}
	}
	r := Result{ShipmentID: shipment, Officer: officer, Status: domain.SecurityPassed, CheckedAt: time.Now().UTC()}
	s.results[shipment] = r
	delete(s.blocked, shipment)
	return r, nil
}
func (s *Service) Get(shipment string) (Result, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.results[shipment]
	if !ok {
		return Result{}, domain.ErrNotFound
	}
	return r, nil
}
func (s *Service) Blocked(shipment string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.blocked[shipment]
	return ok
}
