package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mayankpahwa/aspire-loan-app/app/config"
	"github.com/mayankpahwa/aspire-loan-app/app/routing"
	"github.com/mayankpahwa/aspire-loan-app/internal/repo/mysql"
	"github.com/mayankpahwa/aspire-loan-app/internal/service"
	"github.com/pkg/errors"
)

type serverExec interface {
	ListenAndServe() error
}

// Server represents a generic server
type Server struct {
	Config config.Config
	API    serverExec
}

// New creates a new http server
func New() (*Server, error) {
	conf, err := config.LoadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "could not load envs")
	}
	if err := mysql.InitDatabase(conf); err != nil {
		return nil, errors.Wrap(err, "could not initialize database")
	}

	if err := mysql.LoadSQLFile(conf.MigrationPath); err != nil {
		return nil, errors.Wrap(err, "error running migrations")
	}

	repo := mysql.NewRepo(
		mysql.GetConnection(),
	)
	service := service.NewService(repo)
	handler, err := routing.Handler(conf, service)
	if err != nil {
		return nil, errors.Wrap(err, "could not init handler")
	}
	api := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: handler,
	}
	return &Server{
		Config: conf,
		API:    api,
	}, nil
}

// ListenAndServe will start the server
func (s *Server) ListenAndServe() {
	log.Printf("listening on port: %d", s.Config.Port)
	if err := s.API.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
