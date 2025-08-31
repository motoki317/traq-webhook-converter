package main

import (
	"flag"
	"log/slog"
	"net/http"
	"strconv"
)

var (
	configPath = flag.String("config", "config.yaml", "path to config file")
	config     *Config
)

func main() {
	flag.Parse()

	var err error
	config, err = NewConfig(*configPath)
	if err != nil {
		panic(err)
	}

	tmpl, err := NewTemplater(config.Template)
	if err != nil {
		panic(err)
	}
	s := NewServer(tmpl)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", s.handle)

	slog.Info("Starting server", "version", GetFormattedVersion(), "port", config.Port)
	err = http.ListenAndServe(":"+strconv.Itoa(config.Port), mux)
	if err != nil {
		panic(err)
	}
}
