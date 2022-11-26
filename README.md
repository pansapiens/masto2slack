# masto2slack

`masto2slack` reads posts (toots?) from your Mastodon account and feeds them to a Slack channel.


## Installation

Download the binary from releases, put it somewhere you can run it under a cron job.


## Setup

You'll need to:

- On your Mastodon server, create a 'new application' and get an access token (`access_token`) (you can do this under `https://<your_instance>/settings/applications/new` in the web interface). You'll need `read` access.
- Create a [Slack App](https://api.slack.com/apps?new_app=1) and enable an incoming webhook (you'll get the `webhook_url`)
- Create a configuration file, `~/.config/masto2slack/config.yml`, as below.

Example `~/.config/masto2slack/config.yml`:

```yaml
slack: 
  webhook_url: https://hooks.slack.com/services/TSDFICIU/B0USDFIK9/cs8dyhl3nkj3knfnwnlkfc
mastodon:
  server: https://tootymctootface.social
  access_token: 2903458u2305_34534o8543_34kj5h34k5hk35k3
```

## Running

Either create a cronjob (`crontab -e`) to call `masto2slack` every 5 mins (or longer):

```
*/5 * * * * /usr/local/bin/masto2slack >/dev/null 2>&1
```
(this assumes you've copied `masto2slack` to `/usr/local/bin/masto2slack`)

Alternatively, run:

`watch -n 300 ./masto2slack`

to check for new posts every 5 mins (300 seconds).


## Building from source

You'll need [Go](https://go.dev/dl/) (v1.19+).

```bash
git clone https://github.com/pansapiens/masto2slack.git
cd masto2slack

# Test without building
# (often slow the first time since dependencies will be downloaded)
go run masto2slack.go

# Build and run
go build masto2slack.go
./masto2slack
```

## Ideas / TODO

- Option to run in our own polling loop rather than require a cronjob
- Use websockets so we don't need to poll
- Better logging and a `--verbose` commandline option
- Better Slack post formatting (eg media posts in Slack message 'blocks' ?)