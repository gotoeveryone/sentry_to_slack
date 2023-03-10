package sentry_to_slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor(t *testing.T) {
	for _, p := range []struct {
		level string
		color string
	}{
		{"error", "#ff7738"},
		{"warning", "#b28000"},
		{"info", "#3070e8"},
	} {
		b := Body{Level: p.level}
		assert.Equal(t, b.Color(), p.color)
	}
}
