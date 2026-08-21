package validation

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestInputRules(t *testing.T) {
	if err := Email("ops@example.com"); err != nil {
		t.Fatal(err)
	}
	if err := Email("bad"); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	if err := AirportCode("PEK"); err != nil {
		t.Fatal(err)
	}
	if err := AirportCode("pe"); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	if err := Reference("abc"); err != nil {
		t.Fatal(err)
	}
	if err := Reference("a"); err != domain.ErrInvalid {
		t.Fatal(err)
	}
}
