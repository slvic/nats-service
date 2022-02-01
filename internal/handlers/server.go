package handlers

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/slvic/nats-service/internal/service/deliveries"
)

type Server struct {
	consumer *deliveries.Consumer
}

func NewServer(consumer *deliveries.Consumer) *Server {
	return &Server{consumer: consumer}
}

type Messenger interface {
	Publish(subject string, message []byte) error
	Subscribe(subject string) error
}

func (s *Server) postPublishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// add logger
		return
	}

	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// add logger
		return
	}

	keys, ok := r.URL.Query()["subject"]
	if !ok {
		// add logger
		return
	}

	subject := keys[0]

	err = s.consumer.Publish(subject, message)
	if err != nil {
		// add logger
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (s *Server) postSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// add logger
		return
	}

	keys, ok := r.URL.Query()["subject"]
	if !ok {
		// add logger
		return
	}

	subject := keys[0]

	err := s.consumer.Subscribe(subject)
	if err != nil {
		// add logger
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// add logger
		return
	}

	//
}

func (s *Server) Start() error {
	router := http.NewServeMux()
	router.HandleFunc("/publish", s.postPublishHandler)
	router.HandleFunc("/subscribe", s.postSubscribeHandler)
	address := ":8080"

	srv := http.Server{
		Addr:    address,
		Handler: router,
	}

	done := make(chan os.Signal)
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
