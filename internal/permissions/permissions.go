package permissions

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
)

type Action string

const (
	CreateShipment Action = "shipment:create"
	BookShipment   Action = "shipment:book"
	ReviewCustoms  Action = "customs:review"
	ScanCargo      Action = "cargo:scan"
	ViewSummary    Action = "summary:view"
)

func Allows(role domain.Role, action Action) bool {
	switch role {
	case domain.RoleShipper:
		return action == CreateShipment
	case domain.RoleCoordinator:
		return action == CreateShipment || action == BookShipment || action == ReviewCustoms || action == ViewSummary
	case domain.RoleGroundAgent:
		return action == ScanCargo
	default:
		return false
	}
}
func Require(role domain.Role, action Action) error {
	if !Allows(role, action) {
		return domain.ErrForbidden
	}
	return nil
}
