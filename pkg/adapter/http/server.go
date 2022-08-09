package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/pavel/PSR/pkg/adapter/http/handler"
)

type Server struct {
	router  *chi.Mux
	server  *http.Server
	handler *handler.Handler
}

func New() *Server {
	serv := &Server{}

	serv.router = chi.NewRouter()
	serv.handler = handler.New()

	serv.server = &http.Server{
		Addr:              ":3000",
		Handler:           serv.router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	return serv
}

func (s *Server) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *Server) PrepareRouter() {
	s.router.Get("/", s.handler.GetIndexPage)

	s.router.Get("/game", s.handler.GetGamePage)

	s.router.Get("/echo", s.handler.OpenSocketConnection)

	fs := http.FileServer(http.Dir("./static"))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fs))
}
