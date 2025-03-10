package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"TestKafka/internal/kafka"
	"TestKafka/internal/storage"
)

const (
	kafkaTopic = "messages"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
	db         storage.Creat
	kafka      kafka.Kafka
}

func NewServer(log *slog.Logger, addr string, db storage.Creat, kafka kafka.Kafka) *Server {
	r := chi.NewRouter()
	s := &Server{
		log:   log.With("module", "hendler"),
		db:    db,
		kafka: kafka,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
	r.Handle("/metrics", promhttp.Handler())
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/message", s.TextHandler)
			r.Get("/statistics", s.StatHandler)
		})
	})
	return s
}

func (s *Server) Run() error {
	s.log.Info("server run")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	s.log.Info("server shutdown")
	return nil
}

func (s *Server) TextHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = s.db.CreareNewMessage(r.Context(), req.Message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = s.kafka.Send(kafkaTopic, req.Message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) StatHandler(w http.ResponseWriter, r *http.Request) {

	response, err := http.Get("http://prometheus:9090/api/v1/query?query=kafka_messages_sent_total")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			s.log.Warn("error closing response body")
			return
		}
	}()
	_, err = io.Copy(w, response.Body)
	if err != nil {
		s.log.Warn("error copying response body")
		return
	}
}
