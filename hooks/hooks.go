package hooks

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

// Webhook is an interface that defines all webhooks.
//
// All webhooks should be able to be parsed in order to fill the required fields in each of them.
// They should also have a generalised output for the push notification service to receive.
type Webhook interface {
	parse() error
	body() string
}

// PushEventWebhook is a specialized struct of Webhook that holds all push event information inside other specialized fields
type PushEventWebhook struct {
	repo     string
	branch   string
	pusher   string
	commits  []Commit
	URL      *url.URL
	Original originalWebhook
}

// Parse reads the original webhook and parses it, in order to fill all the necessary fields
// for each struct.
//
// Particularly, this Parse() method fills a PushEventWebhook's fields.
func (w PushEventWebhook) Parse() error {
	var err error

	// get normal entries
	w.repo, err = jsonparser.GetString(w.Original, "repository", "full_name")
	if err != nil {
		w.repo = ""
		return err
	}

	ref, err := jsonparser.GetString(w.Original, "ref")
	if err != nil {
		w.branch = ""
		return err
	}

	refSplit := strings.Split(ref, "/")

	w.branch = refSplit[2]

	w.pusher, err = jsonparser.GetString(w.Original, "pusher", "name")
	if err != nil {
		w.pusher = ""
		return err
	}

	URLString, err := jsonparser.GetString(w.Original, "compare")
	if err != nil {
		return err
	}

	w.URL, err = url.Parse(URLString)
	if err != nil {
		return err
	}

	// iterate through commits object array from original webhook
	_, err = jsonparser.ArrayEach(w.Original, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		// these never fail, webhook object contains these, according to GitHub API documentation
		id, _ := jsonparser.GetString(value, "id")
		committer, _ := jsonparser.GetString(value, "author", "username")
		message, _ := jsonparser.GetString(value, "message")
		URLString, _ := jsonparser.GetString(value, "url")

		URL, err1 := url.Parse(URLString)
		if err1 != nil {
			return
		}

		w.commits = append(w.commits, Commit{
			ID:        id,
			committer: committer,
			message:   message,
			URL:       URL,
		})
	}, "commits")

	// reset entire slice if it fails
	if err != nil {
		w.commits = nil
		return err
	}

	return nil
}

// Body displays the body of a POST request to be sent to a different system.
//
// Currently, GitHookParser supports only Pushover, so the body will be modeled after their API.
func (w PushEventWebhook) Body(apiToken string, userKey string) ([]byte, error) {
	if w.repo == "" || w.pusher == "" || w.commits == nil || w.URL == nil {
		return nil, errors.New("Cannot display body for webhook: Empty values present. Did you call Parse() on the webhook first?")
	}

	var tempMap map[string]interface{}

	if len(w.commits) != 1 {
		tempMap["title"] = "[" + w.repo + ":" + w.branch + "] " + strconv.Itoa(len(w.commits)) + " new commits"
	} else {
		tempMap["title"] = "[" + w.repo + ":" + w.branch + "] " + strconv.Itoa(len(w.commits)) + " new commit"
	}

	var commitString string
	for _, commit := range w.commits {
		// this looks like: a8af291 - Some commit message - commitAuthor
		// this should be clickable, if needed, hence the "a" tag
		commitString += `<a href="` + commit.URL.String() + `">` + commit.ID[len(commit.ID)-7:] + "</a> - " + commit.message + " - " + commit.committer + "\n"
	}

	tempMap["message"] = commitString
	tempMap["url_title"] = "Compare changes"
	tempMap["url"] = w.URL.String()
	tempMap["token"] = apiToken
	tempMap["user"] = userKey
	tempMap["html"] = 1

	body, err := json.Marshal(tempMap)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// IssueEventWebhook is a specialized struct of Webhook that holds all issue event information inside other specialized fields
type IssueEventWebhook struct {
	action   string
	issue    Issue
	URL      *url.URL
	Original originalWebhook
}

// PullRequestEventWebhook is a specialized struct of Webhook that holds all pull request event information inside other specialized fields
type PullRequestEventWebhook struct {
	action      string
	pullRequest PullRequest
	URL         *url.URL
	Original    originalWebhook
}

// Webhook is a struct that holds a generic received webhook
type originalWebhook []byte

// Commit is a struct that represents a commit.
//
// This will be the information presented by a commit to another service.
type Commit struct {
	ID        string
	committer string
	message   string
	URL       *url.URL
}

// Issue is a struct that represents an Issue.
//
// This will be the information presented by an issue to another service.
type Issue struct {
	repo        string
	author      string
	title       string
	description string
	URL         *url.URL
}

// PullRequest is a struct that represents a Pull Request.
//
// This will be the information presented by a pull request to another service.
type PullRequest struct {
	repo        string
	author      string
	title       string
	description string
	URL         *url.URL
}
