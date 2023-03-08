package sentry_to_slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func SendToSlack(r Request) (*http.Response, error) {
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

func CreateRequest(b Body) Request {
	req := Request{
		Channel: os.Getenv("SLACK_CHANNEL"),
		Attachments: []RequestAttachment{
			{
				Title: fmt.Sprintf("<%s|%s>",
					b.Url,
					b.Event.Title,
				),
				Color: b.Color(),
				Fields: []RequestField{
					{Title: "", Value: b.Event.Culprit},
				},
			},
		},
	}

	tags := strings.Split(os.Getenv("TAGS"), ",")
	for _, tag := range tags {
		for _, t := range b.Event.Tags {
			if len(t) < 2 {
				continue
			}
			key := t[0]
			if key != strings.TrimSpace(tag) {
				continue
			}
			value := t[1]
			req.Attachments[0].Fields = append(req.Attachments[0].Fields, RequestField{Title: key, Value: value, Short: true})
		}
	}

	return req
}
