package main

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

type SentryEvent struct {
	Title       string     `json:"title"`
	Culprit     string     `json:"culprit"`
	Environment string     `json:"environment"`
	Tags        [][]string `json:"tags"`
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
