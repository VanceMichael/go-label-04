package postgres

import (
	"context"
	"encoding/json"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"time"
)

func (s *Store) Enqueue(ctx context.Context, v domain.OutboxEvent) error {
	_, err := s.db.Exec(ctx, `INSERT INTO outbox_events(id,tenant_id,topic,aggregate_id,payload,attempts,available_at,last_error) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, v.ID, v.TenantID, v.Topic, v.AggregateID, v.Payload, v.Attempts, v.AvailableAt, v.LastError)
	return err
}
func (s *Store) Claim(ctx context.Context, now time.Time, n int) ([]domain.OutboxEvent, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `WITH picked AS (SELECT id FROM outbox_events WHERE published_at IS NULL AND claimed_at IS NULL AND available_at<= $1 ORDER BY available_at,id FOR UPDATE SKIP LOCKED LIMIT $2) UPDATE outbox_events o SET claimed_at=$1,attempts=o.attempts+1 FROM picked WHERE o.id=picked.id RETURNING o.id,o.tenant_id,o.topic,o.aggregate_id,o.payload,o.attempts,o.available_at,o.claimed_at,o.published_at,o.last_error`, now, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.OutboxEvent, 0)
	for rows.Next() {
		var v domain.OutboxEvent
		var payload []byte
		if err := rows.Scan(&v.ID, &v.TenantID, &v.Topic, &v.AggregateID, &payload, &v.Attempts, &v.AvailableAt, &v.ClaimedAt, &v.PublishedAt, &v.LastError); err != nil {
			return nil, err
		}
		v.Payload = append([]byte(nil), payload...)
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return out, nil
}
func (s *Store) MarkPublished(ctx context.Context, id string, at time.Time) error {
	_, err := s.db.Exec(ctx, `UPDATE outbox_events SET published_at=$2,claimed_at=NULL WHERE id=$1`, id, at)
	return err
}
func (s *Store) MarkFailed(ctx context.Context, id string, at time.Time, msg string) error {
	_, err := s.db.Exec(ctx, `UPDATE outbox_events SET last_error=$2,claimed_at=NULL,available_at=$3 WHERE id=$1`, id, msg, at.Add(time.Minute))
	return err
}
func (s *Store) Summary(ctx context.Context, tenant string) (int, int, error) {
	var p, f int
	err := s.db.QueryRow(ctx, `SELECT count(*) FILTER(WHERE published_at IS NULL),count(*) FILTER(WHERE published_at IS NULL AND attempts>=5) FROM outbox_events WHERE tenant_id=$1`, tenant).Scan(&p, &f)
	return p, f, err
}

var _ = json.Valid
