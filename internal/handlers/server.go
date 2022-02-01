package handlers

import (
	"context"
	"github.com/slvic/nats-service/internal/service/deliveries"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	stream *deliveries.Stream
}

func NewServer(stream *deliveries.Stream) *Server {
	return &Server{stream: stream}
}

type Messenger interface {
	Publish(subject string, message []byte) error
	Subscribe(subject string) error
}

func (s *Server) publishHandler(w http.ResponseWriter, r *http.Request) {
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

	err = s.stream.Publish(subject, message)
	if err != nil {
		// add logger
		return
	}

}
func (s *Server) subscribeHandler(w http.ResponseWriter, r *http.Request) {
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

	err := s.stream.Subscribe(subject)
	if err != nil {
		// add logger
		return
	}
}

func (s *Server) Start() error {
	router := http.NewServeMux()
	router.HandleFunc("/publish", s.publishHandler)
	router.HandleFunc("/subscribe", s.subscribeHandler)
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
