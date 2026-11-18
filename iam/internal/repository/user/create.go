package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ivmerk/goshop/iam/internal/model"
	repoConverter "github.com/ivmerk/goshop/iam/internal/repository/converter"
)

const pgUniqueViolation = "23505"

func (r *repository) Create(ctx context.Context, user *model.User) (string, error) {
	u := repoConverter.ToRepositoryUser(user)

	var uuid string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (login, email, password_hash, notification_methods)
		VALUES ($1, $2, $3, $4)
		RETURNING uuid`,
		u.Login, u.Email, u.PasswordHash, u.NotificationMethods).Scan(&uuid)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return "", model.ErrUserAlreadyExists
		}
		return "", err
	}

	return uuid, nil
}
