package main

import (
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	room "github.com/pavel/PSR/pkg/game"
	web "github.com/pavel/PSR/pkg/room"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	r := chi.NewRouter()

	wRoom := web.NewRoom(&room.GameConfig{
		Name:           "test",
		StepTimeout:    5 * time.Second,
		MaxPlayerCount: 3,
		MaxScore:       7,
		OnlyComputer:   false,
	})
	go wRoom.Main()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse \"templates/index.html\"")
		}

		if err := tmpl.Execute(w, nil); err != nil {
			log.Error().Err(err).Msg("Failed to execute \"templates/index.html\"")
		}
	})

	r.Get("/game", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		tmpl, err := template.ParseFiles("templates/game.html")
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse \"templates/game.html\"")
		}

		err = tmpl.Execute(w, struct {
			ID string
		}{
			ID: id,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to execute \"templates/game.html\"")
		}
	})

	r.Get("/echo", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Upgrade socket error", 500)
			return
		}

		wRoom.AddPlayer(id, conn)
	})

	workDir, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("[server] Can't list work dir")
		return
	}
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/static", filesDir)

	log.Info().Msg("Server started on port 3000")
	http.ListenAndServe(":3000", r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		log.Error().Msg("[server] FileServer does not permit any URL parameters")
		return
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
