package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pavel/PSR/pkg/adapter/http"
	userdb "github.com/pavel/PSR/pkg/adapter/user-db"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	db := userdb.New()
	server := http.New()

	db.Connect(context.Background(), "postgres://postgres:123456@localhost:5432/psr")
	defer db.Disconnect()

	server.PrepareRouter()

	server.Run()
}
