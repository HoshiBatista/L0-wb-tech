package server

import (
	"context"
	"log/slog"
	"net/http"

	"l0-wb-tech/internal/handlers"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(port string, handler *handlers.Handler, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/order/", handler.GetOrderByUID)

	fileServer := http.FileServer(http.Dir("./web"))

	mux.Handle("/", fileServer)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	s.logger.Info("HTTP-сервер запускается", slog.String("port", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("HTTP-сервер останавливается...")
	return s.httpServer.Shutdown(ctx)
}
