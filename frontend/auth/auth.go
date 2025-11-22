package auth

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/run-llama/study-llama/frontend/authdb"
)

var ErrUnauthorized = errors.New("unauthorized")

func AuthorizePost(c *fiber.Ctx) (*db.User, error) {
	sqlDb, err := CreateNewDb()
	if err != nil {
		return nil, ErrUnauthorized
	}
	st := c.Cookies("session_token", "")
	if st == "" {
		return nil, ErrUnauthorized
	}
	queries := db.New(sqlDb)
	ctx := context.Background()
	user, err := queries.GetUserBySessionToken(ctx, pgtype.Text{String: st, Valid: true})
	if err != nil {
		return nil, ErrUnauthorized
	}
	csrf := c.Cookies("csrf_token", "")
	if csrf == "" {
		return nil, ErrUnauthorized
	}
	if csrf != user.CsrfToken.String {
		return nil, ErrUnauthorized
	}
	return &user, nil
}

func AuthorizeGet(c *fiber.Ctx) (*db.User, error) {
	sqlDb, err := CreateNewDb()
	if err != nil {
		return nil, ErrUnauthorized
	}
	st := c.Cookies("session_token", "")
	if st == "" {
		return nil, ErrUnauthorized
	}
	queries := db.New(sqlDb)
	ctx := context.Background()
	user, err := queries.GetUserBySessionToken(ctx, pgtype.Text{String: st, Valid: true})
	if err != nil {
		return nil, ErrUnauthorized
	}
	return &user, nil
}
