package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"time"
)

type Snapshot struct {
	Version   int64
	TenantID  string
	CreatedAt time.Time
	Shipments []domain.Shipment
	Legs      []domain.FlightLeg
	Hash      string
}

func Make(version int64, tenant string, shipments []domain.Shipment, legs []domain.FlightLeg, at time.Time) (Snapshot, error) {
	if version < 1 || tenant == "" || at.IsZero() {
		return Snapshot{}, domain.ErrInvalid
	}
	s := Snapshot{Version: version, TenantID: tenant, CreatedAt: at, Shipments: append([]domain.Shipment(nil), shipments...), Legs: append([]domain.FlightLeg(nil), legs...)}
	sort.Slice(s.Shipments, func(i, j int) bool { return s.Shipments[i].ID < s.Shipments[j].ID })
	sort.Slice(s.Legs, func(i, j int) bool { return s.Legs[i].ID < s.Legs[j].ID })
	s.Hash = hash(s)
	return s, nil
}
func (s Snapshot) EqualPayload(other Snapshot) bool {
	return s.Hash == other.Hash && s.TenantID == other.TenantID && s.Version == other.Version
}
func (s Snapshot) Validate() error {
	if s.Version < 1 || s.TenantID == "" || s.Hash == "" {
		return domain.ErrInvalid
	}
	if s.Hash != hash(s) {
		return domain.ErrConflict
	}
	return nil
}
func hash(s Snapshot) string {
	b, _ := json.Marshal(struct {
		V int64
		T string
		S []domain.Shipment
		L []domain.FlightLeg
	}{s.Version, s.TenantID, s.Shipments, s.Legs})
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
