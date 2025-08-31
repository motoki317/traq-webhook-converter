package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type m = map[string]any

type Server struct {
	templater *Templater
}

func NewServer(templater *Templater) *Server {
	return &Server{templater: templater}
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	// Parse
	var reqData m
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		slog.Error("parsing request body", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Template
	msg, err := s.templater.template(reqData)
	if err != nil {
		slog.Error("templating", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Post
	err = PostWebhook(config.Webhook.URL, msg, config.Webhook.Secret)
	if err != nil {
		slog.Error("posting webhook", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
