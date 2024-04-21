package pg

import (
	"context"
	"fmt"

	api "github.com/eriktate/splitdecision"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (pg PG) CreateUser(ctx context.Context, user api.NewUser) (api.ID, error) {
	id := api.NewID()

	query := `
	insert into users (id, email, auth_type, token)
	values ($1, $2, $3, $4)
	`

	if _, err := pg.pool.Exec(ctx, query, id, user.Email, user.AuthType, user.Token); err != nil {
		return api.NullID(), fmt.Errorf("failed to CreateUser: %w", err)
	}

	return id, nil
}

func (pg PG) GetUser(ctx context.Context, id api.ID) (api.User, error) {
	var user api.User
	query := "select * from users where id = $1"

	if err := pgxscan.Get(ctx, pg.pool, &user, query, id); err != nil {
		return user, fmt.Errorf("failed to GetUser: %w", err)
	}

	return user, nil
}

func (pg PG) GetUserByEmail(ctx context.Context, email string) (api.User, error) {
	var user api.User
	query := "select * from users where email = $1"

	if err := pgxscan.Get(ctx, pg.pool, &user, query, email); err != nil {
		return user, fmt.Errorf("failed to GetUserByEmail: %w", err)
	}

	return user, nil
}

func (pg PG) DeleteUser(ctx context.Context, id api.ID) error {
	query := "delete from users where id = $1"

	if _, err := pg.pool.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to DeleteUser: %w", err)
	}

	return nil
}
