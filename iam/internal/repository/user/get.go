package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/ivmerk/goshop/iam/internal/model"
	repoConverter "github.com/ivmerk/goshop/iam/internal/repository/converter"
	repoModel "github.com/ivmerk/goshop/iam/internal/repository/model"
)

const selectUser = `SELECT uuid, login, email, password_hash, notification_methods, created_at, updated_at FROM users`

func (r *repository) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	return r.get(ctx, selectUser+` WHERE uuid = $1`, uuid)
}

func (r *repository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	return r.get(ctx, selectUser+` WHERE login = $1`, login)
}

func (r *repository) get(ctx context.Context, query string, arg any) (*model.User, error) {
	var u repoModel.User
	row := r.pool.QueryRow(ctx, query, arg)

	err := row.Scan(&u.UUID, &u.Login, &u.Email, &u.PasswordHash, &u.NotificationMethods, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	return repoConverter.FromRepositoryUser(&u), nil
}
