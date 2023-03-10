package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	app "github.com/gotoeveryone/sentry_to_slack/src"
)

func sendToSlack(r app.Request) (*http.Response, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	c := http.Client{}
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_API_TOKEN")))
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func createRequest(e app.Event) app.Request {
	req := app.Request{
		Channel: os.Getenv("SLACK_CHANNEL"),
		Attachments: []app.RequestAttachment{
			{
				Title: fmt.Sprintf("<%s|%s>",
					e.Url,
					e.Event.Title,
				),
				Color: e.Color(),
				Fields: []app.RequestField{
					{Title: "", Value: e.Event.Culprit},
				},
			},
		},
	}

	tags := strings.Split(os.Getenv("TAGS"), ",")
	for _, tag := range tags {
		for _, t := range e.Event.Tags {
			if len(t) < 2 {
				continue
			}
			key := t[0]
			if key != strings.TrimSpace(tag) {
				continue
			}
			value := t[1]
			req.Attachments[0].Fields = append(req.Attachments[0].Fields, app.RequestField{Title: key, Value: value, Short: true})
		}
	}

	return req
}

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (string, error) {
	var e app.Event
	if err := json.Unmarshal([]byte(request.Body), &e); err != nil {
		return "error", err
	}
	log.Printf("An error has occurred [%s]: %s", e.Event.Environment, e.Event.Title)
	req := createRequest(e)
	res, err := sendToSlack(req)
	if err != nil {
		return "error", err
	}

	return fmt.Sprintf("success: %d", res.StatusCode), nil
}

func main() {
	if os.Getenv("DEBUG") == "1" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				w.WriteHeader(405)
				fmt.Fprint(w, "Method not allowed")
				return
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprint(w, err.Error())
				return
			}
			if _, err := HandleRequest(context.Background(), events.LambdaFunctionURLRequest{Body: string(body)}); err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, err.Error())
				return
			}
			w.WriteHeader(204)
		})
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalln(err)
		}
	}
	lambda.Start(HandleRequest)
}
