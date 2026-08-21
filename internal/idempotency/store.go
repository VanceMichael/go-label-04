package idempotency

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
	"time"
)

type Record struct {
	TenantID    string
	Key         string
	Fingerprint string
	Status      int
	Body        []byte
	CreatedAt   time.Time
	ExpiresAt   time.Time
}
type Store struct {
	mu      sync.Mutex
	records map[string]Record
}

func New() *Store { return &Store{records: map[string]Record{}} }
func Fingerprint(method, path string, body []byte) string {
	sum := sha256.Sum256(append([]byte(method+"\n"+path+"\n"), body...))
	return hex.EncodeToString(sum[:])
}
func (s *Store) Begin(tenant, key, fingerprint string, now time.Time) (Record, bool, error) {
	if tenant == "" || key == "" || fingerprint == "" {
		return Record{}, false, domain.ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	compound := tenant + "|" + key
	if old, ok := s.records[compound]; ok && now.Before(old.ExpiresAt) {
		if old.Fingerprint != fingerprint {
			return Record{}, false, domain.ErrConflict
		}
		return old, true, nil
	}
	v := Record{TenantID: tenant, Key: key, Fingerprint: fingerprint, CreatedAt: now, ExpiresAt: now.Add(24 * time.Hour), Status: 102}
	s.records[compound] = v
	return v, false, nil
}
func (s *Store) Complete(tenant, key string, status int, body []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	compound := tenant + "|" + key
	v, ok := s.records[compound]
	if !ok {
		return domain.ErrNotFound
	}
	v.Status = status
	v.Body = append([]byte(nil), body...)
	s.records[compound] = v
	return nil
}
func (s *Store) Cleanup(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for k, v := range s.records {
		if !now.Before(v.ExpiresAt) {
			delete(s.records, k)
			n++
		}
	}
	return n
}
