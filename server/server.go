package server

import (
	"database/sql"
	"net/http"

	"github.com/stormfirefox1/GitHookParser/config"
)

// Server is the struct that will contain all the information the standard library needs in order to run correctly.
type Server struct {
	Router *http.ServeMux
	Env    config.Env
	DB     *sql.DB
}

// Routes configures all of the functions for the whole server
func (s *Server) Routes() {
	s.Router = http.NewServeMux()
	s.Router.HandleFunc("/github-hook", s.handleGitHubHook())
}
