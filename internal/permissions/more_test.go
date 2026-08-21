package permissions

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestRequire(t *testing.T) {
	if err := Require(domain.RoleCoordinator, BookShipment); err != nil {
		t.Fatal(err)
	}
	if err := Require(domain.RoleGroundAgent, BookShipment); err != domain.ErrForbidden {
		t.Fatal(err)
	}
}
