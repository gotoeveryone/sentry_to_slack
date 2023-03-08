package sentry_to_slack

type Body struct {
	Url     string `json:"url"`
	Project string `json:"project"`
	Level   string `json:"level"`
	Event   Event  `json:"event"`
}

func (b *Body) Color() string {
	switch b.Level {
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

type Event struct {
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
