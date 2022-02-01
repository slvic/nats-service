package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/slvic/nats-service/internal/types"
	"go.uber.org/zap"
)

type StoreService interface {
	GetMessages(string) types.Order
}

type Server struct {
	storeService StoreService
	logger       *zap.Logger // should probably pass logger as an argument
}

func NewServer(storeService StoreService, logger *zap.Logger) *Server {
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

	keys, found := r.URL.Query()["subject"]
	if !found {
		notFoundResponse, err := json.Marshal(map[string]string{"reason": "there is no such a subject"})
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
	subject := keys[0]
	s.storeService.GetMessages(subject)
}

func (s *Server) Start() error {
	router := http.NewServeMux()
	router.HandleFunc("/messages", s.getMessagesHandler)
	address := ":8080"

	srv := http.Server{
		Addr:    address,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("listen error: %v", err)
		}
	}()
	log.Println("Server started")

	<-done
	log.Println("Server stopped")

	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
		return err
	}
	log.Println("Server exited properly")
	return nil
}
