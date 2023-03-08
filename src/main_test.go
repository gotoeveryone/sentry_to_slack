package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRequest(t *testing.T) {
	os.Setenv("SLACK_CHANNEL", "test_channel")
	os.Setenv("TAGS", "server_name ,environment")
	e := Event{
		Url:     "https://test.example.com",
		Level:   "error",
		Project: "test_project",
		Event: SentryEvent{
			Title:   "TestError",
			Culprit: "This is test error",
			Tags: [][]string{
				{"environment", "test"},
				{"server_name", "test.example.com"},
				{"ignore", "test"},
			},
		},
	}
	r := createRequest(e)
	assert.Equal(t, r, Request{
		Channel: "test_channel",
		Attachments: []RequestAttachment{
			{
				Title: fmt.Sprintf("<%s|%s>",
					e.Url,
					e.Event.Title,
				),
				Color: e.Color(),
				Fields: []RequestField{
					{Title: "", Value: e.Event.Culprit},
					{Title: "server_name", Value: "test.example.com", Short: true},
					{Title: "environment", Value: "test", Short: true},
				},
			},
		},
	})
}
