package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"time"
)

type Item struct {
	SKU         string
	Description string
	Quantity    int
	WeightKg    int64
	Hazardous   bool
}
type Manifest struct {
	ID         string
	ShipmentID string
	Version    int
	Items      []Item
	SealedAt   *time.Time
	Hash       string
}

func Validate(items []Item) error {
	if len(items) == 0 {
		return domain.ErrInvalid
	}
	for _, item := range items {
		if item.SKU == "" || item.Quantity <= 0 || item.WeightKg <= 0 {
			return domain.ErrInvalid
		}
	}
	return nil
}
func Build(id, shipment string, items []Item) (Manifest, error) {
	if err := Validate(items); err != nil {
		return Manifest{}, err
	}
	copyItems := append([]Item(nil), items...)
	sort.Slice(copyItems, func(i, j int) bool { return copyItems[i].SKU < copyItems[j].SKU })
	m := Manifest{ID: id, ShipmentID: shipment, Version: 1, Items: copyItems}
	m.Hash = hash(m)
	return m, nil
}
func (m Manifest) Seal(now time.Time) (Manifest, error) {
	if m.SealedAt != nil {
		return Manifest{}, domain.ErrConflict
	}
	m.SealedAt = &now
	m.Hash = hash(m)
	return m, nil
}
func (m Manifest) TotalWeight() int64 {
	var n int64
	for _, i := range m.Items {
		n += i.WeightKg * int64(i.Quantity)
	}
	return n
}
func (m Manifest) HazardousCount() int {
	n := 0
	for _, i := range m.Items {
		if i.Hazardous {
			n += i.Quantity
		}
	}
	return n
}
func hash(m Manifest) string {
	b, _ := json.Marshal(struct {
		ID       string
		Shipment string
		Version  int
		Items    []Item
	}{m.ID, m.ShipmentID, m.Version, m.Items})
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
