package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"

	"github.com/stormfirefox1/GitHookParser/hooks"
	"github.com/stormfirefox1/GitHookParser/log"
)

func (s *server) handleGitHubHook() http.HandlerFunc {
	var (
		init sync.Once
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			s.startDB()
		})

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Bad method! Use POST")
			return
		}

		eventType := r.Header.Get("X-GitHub-Event")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Can't read body")
			log.Info(logrus.Fields{
				"timestamp":    time.Now(),
				"Host":         r.Host,
				"RemoteAddr":   r.RemoteAddr,
				"Method":       r.Method,
				"UserAgent":    r.Header.Get("User-Agent"),
				"EventType":    eventType,
				"ResponseCode": http.StatusInternalServerError,
				"error":        err,
			}, "handleGitHook hit")
			return
		}

		switch eventType {
		case "push":
			webhook := hooks.PushEventWebhook{
				original: body,
			}
			err = webhook.parse()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Can't parse webhook")
				log.Info(logrus.Fields{
					"timestamp":    time.Now(),
					"Host":         r.Host,
					"RemoteAddr":   r.RemoteAddr,
					"Method":       r.Method,
					"UserAgent":    r.Header.Get("User-Agent"),
					"EventType":    eventType,
					"ResponseCode": http.StatusInternalServerError,
					"error":        err,
				}, "handleGitHook hit")
				return
			}
			hookBody, err := webhook.body()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Can't parse webhook")
				log.Info(logrus.Fields{
					"timestamp":    time.Now(),
					"Host":         r.Host,
					"RemoteAddr":   r.RemoteAddr,
					"Method":       r.Method,
					"UserAgent":    r.Header.Get("User-Agent"),
					"EventType":    eventType,
					"ResponseCode": http.StatusInternalServerError,
					"error":        err,
				}, "handleGitHook hit")
				return
			}
			hookBodyReader := bytes.NewReader(hookBody)
			resp, err := http.Post(s.env.Get("REDIRECT_URL").(string), "application/json", hookBodyReader)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Can't redirect webhook")
				log.Info(logrus.Fields{
					"timestamp":    time.Now(),
					"Host":         r.Host,
					"RemoteAddr":   r.RemoteAddr,
					"Method":       r.Method,
					"UserAgent":    r.Header.Get("User-Agent"),
					"EventType":    eventType,
					"ResponseCode": http.StatusInternalServerError,
					"error":        err,
				}, "handleGitHook hit")
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Success")
			log.Info(logrus.Fields{
				"timestamp":    time.Now(),
				"Host":         r.Host,
				"RemoteAddr":   r.RemoteAddr,
				"Method":       r.Method,
				"UserAgent":    r.Header.Get("User-Agent"),
				"EventType":    eventType,
				"ResponseCode": http.StatusOK,
			}, "handleGitHook hit")
			err = s.addHook(body, "github")
			if err != nil {
				log.Info(logrus.Fields{
					"timestamp": time.Now(),
					"error":     err,
					"category":  "github",
				}, "handleGitHook save failure - couldn't save hook")
			}
			return
		case "ping":
			zen, _ := jsonparser.GetString(body, "zen")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, zen)
			log.Info(logrus.Fields{
				"timestamp":  time.Now(),
				"Host":       r.Host,
				"RemoteAddr": r.RemoteAddr,
				"Method":     r.Method,
				"UserAgent":  r.Header.Get("User-Agent"),
				"EventType":  eventType,
				"Zen":        zen,
			}, "handleGitHook ping hit")
			err = s.addHook(body, "github")
			if err != nil {
				log.Info(logrus.Fields{
					"timestamp": time.Now(),
					"error":     err,
					"category":  "github",
				}, "handleGitHook save failure - couldn't save hook")
			}
			return
		default:
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprint(w, "Event type not supported by handler")
			log.Info(logrus.Fields{
				"timestamp":    time.Now(),
				"Host":         r.Host,
				"RemoteAddr":   r.RemoteAddr,
				"Method":       r.Method,
				"UserAgent":    r.Header.Get("User-Agent"),
				"EventType":    eventType,
				"ResponseCode": http.StatusNotImplemented,
			}, "handleGitHook hit")
			return
		}
	}
}
