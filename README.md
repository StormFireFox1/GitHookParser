# GitHookParser

GitHookParser is a nifty server written in Go that can take a GitHub webhook, parse it according to configuration, and redirect it to another POST endpoint.

Primarily, this can be used to manage multiple webhooks for GitHub repos on the fly and send them all to a push notification service defined by the user,
such as [Pushover](https://pushover.net/).