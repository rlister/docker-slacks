package main

import (
	"bytes"
	"fmt"
	"github.com/samalba/dockerclient"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
)

// data to use in template
type Data struct {
	Timestamp time.Time
	Hostname  string
	Event     *dockerclient.Event
}

// error handler
func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}

func waitForInterrupt() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	for _ = range sigChan {
		os.Exit(0)
	}
}

// return closure for sending to given slack webhook
func setupPostToSlack(webhook string) func(string) {
	return func(msg string) {
		resp, err := http.PostForm(webhook, url.Values{"payload": {msg}})
		check(err)

		if resp.StatusCode != http.StatusOK {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			check(err)
			fmt.Printf("error posting to slack: [%s] %s\n", resp.Status, body)
		}
	}
}

// return closure for handling docker events
func setupCallback(filename string) func(*dockerclient.Event, chan error, ...interface{}) {

	// hostname to report
	hostname, err := os.Hostname()
	check(err)

	// function to send data to slack
	postToSlack := setupPostToSlack(os.Getenv("WEBHOOK"))

	// load a template from file
	tmpl, err := template.ParseFiles(filename)
	check(err)

	// format json from template and send to slack
	return func(event *dockerclient.Event, ec chan error, args ...interface{}) {
		fmt.Printf("%+v\n", *event) // log to stdout

		data := Data{
			Timestamp: time.Unix(event.Time, 0),
			Hostname:  hostname,
			Event:     event,
		}

		var output bytes.Buffer      // will contain output to send
		tmpl.Execute(&output, &data) // format json from template
		postToSlack(output.String()) // send json
	}
}

func main() {

	// set file to load json template
	template_file := os.Getenv("TEMPLATE_FILE")
	if template_file == "" {
		template_file = "default.json"
	}

	// can listen on given http url, or default to unix socket
	host := os.Getenv("DOCKER_HOST")
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}

	// docker client
	docker, err := dockerclient.NewDockerClient(host, nil)
	check(err)

	// listen to events
	callback := setupCallback(template_file)
	docker.StartMonitorEvents(callback, nil)

	// wait forever
	waitForInterrupt()
}
