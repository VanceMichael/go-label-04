package permissions

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestRoleMatrix(t *testing.T) {
	if !Allows(domain.RoleShipper, CreateShipment) {
		t.Fatal("shipper create")
	}
	if Allows(domain.RoleShipper, BookShipment) {
		t.Fatal("shipper book")
	}
	if !Allows(domain.RoleCoordinator, ViewSummary) {
		t.Fatal("coordinator summary")
	}
	if !Allows(domain.RoleGroundAgent, ScanCargo) {
		t.Fatal("ground scan")
	}
}
