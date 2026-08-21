package audit

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/VanceMichael/go-base-airbridge/internal/repository"
)

type Service struct{ repo repository.AuditRepository }

func New(repo repository.AuditRepository) *Service { return &Service{repo: repo} }
func (s *Service) Record(ctx context.Context, e domain.AuditEvent) error {
	if e.ID == "" || e.TenantID == "" || e.ActorID == "" || e.ObjectType == "" || e.ObjectID == "" || e.Action == "" {
		return domain.ErrInvalid
	}
	return s.repo.AppendAudit(ctx, e)
}
func (s *Service) Recent(ctx context.Context, tenant string, req domain.PageRequest) (domain.Page[domain.AuditEvent], error) {
	return s.repo.ListAudit(ctx, tenant, req)
}
