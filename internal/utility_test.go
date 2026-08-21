package internal_test

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/collections"
	"github.com/VanceMichael/go-base-airbridge/internal/expiry"
	"github.com/VanceMichael/go-base-airbridge/internal/geo"
	"github.com/VanceMichael/go-base-airbridge/internal/locale"
	"github.com/VanceMichael/go-base-airbridge/internal/ratelimit"
	"github.com/VanceMichael/go-base-airbridge/internal/ratios"
	"github.com/VanceMichael/go-base-airbridge/internal/sequence"
	"github.com/VanceMichael/go-base-airbridge/internal/sets"
	"github.com/VanceMichael/go-base-airbridge/internal/stringsx"
	"github.com/VanceMichael/go-base-airbridge/internal/version"
	"testing"
	"time"
)

func TestUtilityContracts(t *testing.T) {
	if got := collections.Chunk([]int{1, 2, 3}, 2); len(got) != 2 || len(got[1]) != 1 {
		t.Fatal(got)
	}
	if got := collections.Map([]int{1, 2}, func(v int) int { return v * 2 }); got[1] != 4 {
		t.Fatal(got)
	}
	if got := sets.Union([]string{"a"}, []string{"a", "b"}); len(got) != 2 {
		t.Fatal(got)
	}
	if stringsx.Slug(" Cargo  PEK ") != "cargo-pek" {
		t.Fatal(stringsx.Slug(" Cargo  PEK "))
	}
}
func TestExpiryAndTime(t *testing.T) {
	now := time.Now()
	v, err := expiry.New(now, time.Hour, time.Minute)
	if err != nil || !v.Active(now.Add(time.Minute)) {
		t.Fatal(err)
	}
	if v.Active(now.Add(2 * time.Hour)) {
		t.Fatal("expired value active")
	}
	if _, err := locale.Parse("Asia/Shanghai", "2026-01-01T08:00:00+08:00"); err != nil {
		t.Fatal(err)
	}
	if !locale.BusinessDay("Asia/Shanghai", now) {
		t.Log("weekend is valid result")
	}
}
func TestRatesAndGeo(t *testing.T) {
	p, err := ratios.Percent(25, 100)
	if err != nil || p != 25 {
		t.Fatal(p, err)
	}
	d, err := geo.Distance(geo.Point{Lat: 0, Lon: 0}, geo.Point{Lat: 0, Lon: 1})
	if err != nil || d < 100 {
		t.Fatal(d, err)
	}
}
func TestRateLimiterSequenceVersion(t *testing.T) {
	now := time.Now()
	l := ratelimit.New(2, time.Minute)
	if !l.Allow("k", now) || !l.Allow("k", now) {
		t.Fatal("first two denied")
	}
	if l.Allow("k", now) {
		t.Fatal("third allowed")
	}
	g := sequence.New("AB")
	a, _ := g.Next(now)
	b, _ := g.Next(now)
	if a == b {
		t.Fatal("duplicate sequence")
	}
	v, err := version.Parse("v1.2.3")
	if err != nil || !v.AtLeast(version.Value{Major: 1, Minor: 1, Patch: 9}) {
		t.Fatal(v, err)
	}
	_, cancel := context.WithCancel(context.Background())
	cancel()
}
