package pg

import (
	"context"
	"fmt"
	"time"

	api "github.com/eriktate/splitdecision"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (pg PG) CreateSession(ctx context.Context, session api.NewSession) (api.ID, error) {
	id := api.NewID()

	query := `
	insert into sessions (id, user_id, expires_at)
	values ($1, $2, $3)
	`

	if _, err := pg.pool.Exec(ctx, query, id, session.UserID, session.ExpiresAt); err != nil {
		return api.NullID(), fmt.Errorf("failed to CreateSession: %w", err)
	}

	return id, nil
}

func (pg PG) GetSession(ctx context.Context, id api.ID) (api.Session, error) {
	var session api.Session
	query := "select * from sessions where id = $1 and expires_at > now()"

	if err := pgxscan.Get(ctx, pg.pool, &session, query, id); err != nil {
		return session, fmt.Errorf("failed to GetSession: %w", err)
	}

	return session, nil
}

func (pg PG) GetUserSession(ctx context.Context, userID api.ID) (api.Session, error) {
	var session api.Session
	query := "select * from sessions where user_id = $1 and expires_at > now() order by expires_at desc"

	if err := pgxscan.Get(ctx, pg.pool, &session, query, userID); err != nil {
		return session, fmt.Errorf("failed to GetUserSession: %w", err)
	}

	return session, nil
}

func (pg PG) ExtendSession(ctx context.Context, id api.ID, expiresAt time.Time) error {
	query := "update sessions set expires_at = id where id = $1 and expires_at < now()"

	if _, err := pg.pool.Exec(ctx, query, id, expiresAt); err != nil {
		return fmt.Errorf("failed to ExtendSession: %w", err)
	}

	return nil
}

func (pg PG) DeleteSession(ctx context.Context, id api.ID) error {
	query := "delete sessions where id = $1"

	if _, err := pg.pool.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to DeleteSession: %w", err)
	}

	return nil
}

func (pg PG) DeleteUserSessions(ctx context.Context, userID api.ID) error {
	query := "delete sessions where user_id = $1"

	if _, err := pg.pool.Exec(ctx, query, userID); err != nil {
		return fmt.Errorf("failed to DeleteUserSessions: %w", err)
	}

	return nil
}
