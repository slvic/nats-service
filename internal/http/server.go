package http

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"go.uber.org/zap"
)

type StoreService interface {
	GetMessageById(string) ([]byte, error)
}

type Server struct {
	storeService StoreService
	logger       *zap.Logger
	httpServer   *http.Server
}

func New(storeService StoreService, logger *zap.Logger) *Server {
	return &Server{
		storeService: storeService,
		logger:       logger,
	}
}

func (s *Server) getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.logger.Info("used wrong http method for",
			zap.String("want", http.MethodGet),
			zap.String("got", r.Method))
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	keys, found := r.URL.Query()["id"]
	if !found {
		s.logger.Info("wrong query parameter")
		notFoundResponse, err := json.Marshal(map[string]string{"reason": "wrong query parameter"})
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(notFoundResponse)
		if err != nil {
			return
		}
	}
	id := keys[0]
	order, err := s.storeService.GetMessageById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := http.NewServeMux()
	router.HandleFunc("/messages", s.getMessagesHandler)
	address := ":3000"

	srv := http.Server{
		Addr:        address,
		Handler:     router,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	s.httpServer = &srv

	return srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
