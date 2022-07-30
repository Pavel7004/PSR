package main

import (
	"os"

	"github.com/pavel/PSR/pkg/adapter/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	server := http.New()

	server.PrepareRouter()

	server.Run()
}
