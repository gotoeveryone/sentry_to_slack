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
)

type SentryEventMetadata struct {
	Filename string `json:"filename"`
	Function string `json:"function"`
}

type SentryEvent struct {
	Title       string              `json:"title"`
	Culprit     string              `json:"culprit"`
	Environment string              `json:"environment"`
	Tags        [][]string          `json:"tags"`
	Metadata    SentryEventMetadata `json:"metadata"`
}

type Event struct {
	Url     string      `json:"url"`
	Project string      `json:"project"`
	Level   string      `json:"level"`
	Event   SentryEvent `json:"event"`
}

func (e *Event) Color() string {
	switch e.Level {
	case "error":
		return "#ff7738"
	case "warning":
		return "#b28000"
	case "info":
		return "#3070e8"
	default:
		return ""
	}
}

// Slack へリクエストするための構造体
type Request struct {
	Channel     string              `json:"channel"`
	Attachments []RequestAttachment `json:"attachments"`
}

type RequestAttachment struct {
	Title  string         `json:"title"`
	Color  string         `json:"color"`
	Fields []RequestField `json:"fields"`
}

type RequestField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func sendToSlack(r Request) (*http.Response, error) {
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

func HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (string, error) {
	var e Event
	if err := json.Unmarshal([]byte(request.Body), &e); err != nil {
		return "error", err
	}
	log.Printf("An error has occurred [%s]: %s", e.Event.Environment, e.Event.Title)
	p := Request{
		Channel: os.Getenv("SLACK_CHANNEL"),
		Attachments: []RequestAttachment{
			{
				Title: fmt.Sprintf("<%s|%s>",
					e.Url,
					e.Event.Title,
				),
				Color: e.Color(),
				Fields: []RequestField{
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
			if key != tag {
				continue
			}
			value := t[1]
			p.Attachments[0].Fields = append(p.Attachments[0].Fields, RequestField{Title: key, Value: value, Short: true})
		}
	}

	res, err := sendToSlack(p)
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
