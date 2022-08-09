package userdb

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
)

type UserDB struct {
	conn *pgx.Conn
}

func New() *UserDB {
	return &UserDB{}
}

func (db *UserDB) Connect(ctx context.Context, uri string) error {
	conn, err := pgx.Connect(ctx, uri)
	if err != nil {
		log.Error().Err(err).Msgf("Can't connect to db at %q.", uri)
		return err
	}

	db.conn = conn

	if err := db.conn.Ping(ctx); err != nil {
		log.Error().Err(err).Msgf("Connected to db at %q, but failed to ping it.", uri)
		return err
	}

	log.Info().Msgf("Connected to db %q", uri)
	return nil
}

func (db *UserDB) Disconnect() error {
	return db.conn.Close(context.Background())
}

// func (db *UserDB) AddUser(ctx context.Context, )
