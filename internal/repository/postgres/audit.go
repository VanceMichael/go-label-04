package postgres

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
)

func (s *Store) AppendAudit(ctx context.Context, v domain.AuditEvent) error {
	_, err := s.db.Exec(ctx, `INSERT INTO audit_events(id,tenant_id,actor_id,object_type,object_id,action,result,request_id,occurred_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, v.ID, v.TenantID, v.ActorID, v.ObjectType, v.ObjectID, v.Action, v.Result, v.RequestID, v.OccurredAt)
	return err
}
func (s *Store) ListAudit(ctx context.Context, tenant string, req domain.PageRequest) (domain.Page[domain.AuditEvent], error) {
	req = req.Normalized()
	rows, err := s.db.Query(ctx, `SELECT id,tenant_id,actor_id,object_type,object_id,action,result,request_id,occurred_at FROM audit_events WHERE tenant_id=$1 ORDER BY occurred_at,id LIMIT $2`, tenant, req.Limit)
	if err != nil {
		return domain.Page[domain.AuditEvent]{}, err
	}
	defer rows.Close()
	out := make([]domain.AuditEvent, 0)
	for rows.Next() {
		var v domain.AuditEvent
		if err := rows.Scan(&v.ID, &v.TenantID, &v.ActorID, &v.ObjectType, &v.ObjectID, &v.Action, &v.Result, &v.RequestID, &v.OccurredAt); err != nil {
			return domain.Page[domain.AuditEvent]{}, err
		}
		out = append(out, v)
	}
	return domain.Page[domain.AuditEvent]{Items: out, Total: len(out)}, rows.Err()
}
