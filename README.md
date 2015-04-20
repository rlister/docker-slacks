# docker-slacks

A tiny golang daemon and container to watch your docker daemon and
sends container events to your [Slack](https://slack.com) chat.

This is inspired by [slack-docker](https://github.com/int128/slack-docker).

## Setup

Set up a Slack
[incoming webhook](https://my.slack.com/services/new/incoming-webhook),
and use the URL to set the `WEBHOOK` environment variable.

## Run natively

Compile for your program using go, and run as follows:

```
WEBHOOK=<url> docker-slacks
```

## Run in docker

Pull `rlister/docker-slacks` from docker hub, and run as follows. You
will probably want to pass in your real hostname to report to slack
(instead of the container ID). Mount `docker.sock` so the
container can read docker daemon events.

```
docker run \
  -h $(hostname) \
  -e WEBHOOK=<url> \
  -v /var/run/docker.sock:/var/run/docker.sock
  rlister/docker-slacks
```

## Docker API

## Slack message templates
