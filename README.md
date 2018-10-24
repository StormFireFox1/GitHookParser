# GitHookParser

GitHookParser is a nifty server written in Go that can take a GitHub webhook, parse it according to configuration, and redirect it to another POST endpoint.

Primarily, this can be used to manage multiple webhooks for GitHub repos on the fly and send them all to a push notification service defined by the user,
such as [Pushover](https://pushover.net/). Currently, this is the only notification service supported.

## Getting started

To install GitHookParser, simply

```
go get -u -v github.com/stormfirefox1/GitHookParser
```

By default, a configuration file is needed under *$HOME/.config/git-hook-parser/config.yml*. If there is a file there, it will read any environment variables
from the respective file. You will need to set, at minimum:

- REDIRECT_URL
- API_KEY
- USER_TOKEN

After creating the configuration file, run the compiled binary, either by `go install` it or by running `go build` and moving the binary to your directory of choice.

Then, point GitHub's webhook to:

```
http://example.com:80/github-hook
```

and watch the magic happen.