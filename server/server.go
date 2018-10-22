package server

import (
	"database/sql"
	"net/http"

	"github.com/stormfirefox1/GitHookParser/config"
)

type server struct {
	router *http.ServeMux
	env    config.Env
	db     *sql.DB
}

// routes configures all of the functions for the whole server
func (s *server) routes() {
	s.router.HandleFunc("/github-hook", s.handleGitHubHook())
}
