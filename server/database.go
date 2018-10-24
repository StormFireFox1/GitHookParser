package server

import (
	"database/sql"
	"fmt"
	"time"

	// this is used to import the SQLite3 driver used by the project.
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/stormfirefox1/GitHookParser/log"
)

// startDB initializes a database in a separate directory.
//
// This should only be run once.
func (s *Server) startDB() {
	db, err := sql.Open("sqlite3", "./db/hooks.db")
	if err != nil {
		log.Fatal(logrus.Fields{}, fmt.Errorf("Error starting up database: %v", err))
	}
	s.DB = db

}

// addHook adds a webhook to the database in the correct category.
//
// An error is returned if the database throws an error.
func (s *Server) addHook(hook []byte, category string) error {
	statement, err := s.DB.Prepare("INSERT INTO hooks(hook, category, created) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(hook, category, time.Now().String())
	if err != nil {
		return err
	}

	return nil
}
