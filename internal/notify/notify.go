// Package notify provides notification dispatch for depwatch alerts.
// It supports writing structured alert summaries to configurable output
// channels such as stdout or a webhook endpoint.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/depwatch/internal/alert"
	"github.com/depwatch/internal/summary"
)

// Channel represents a notification destination.
type Channel string

const (
	ChannelStdout  Channel = "stdout"
	ChannelWebhook Channel = "webhook"
)

// Config holds notification settings.
type Config struct {
	Channel    Channel
	WebhookURL string
}

// Payload is the JSON body sent to a webhook.
type Payload struct {
	Timestamp    time.Time      `json:"timestamp"`
	OverallLevel string         `json:"overall_level"`
	Total        int            `json:"total"`
	Critical     int            `json:"critical"`
	Warn         int            `json:"warn"`
	Ecosystems   []string       `json:"ecosystems"`
	Alerts       []alert.Alert  `json:"alerts"`
}

// Notifier dispatches notifications.
type Notifier struct {
	cfg    Config
	out    io.Writer
	client *http.Client
}

// New returns a Notifier using the given Config.
func New(cfg Config) *Notifier {
	return &Notifier{
		cfg:    cfg,
		out:    os.Stdout,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// NewWithWriter returns a Notifier that writes stdout notifications to w.
func NewWithWriter(cfg Config, w io.Writer) *Notifier {
	n := New(cfg)
	n.out = w
	return n
}

// Send dispatches alerts via the configured channel.
func (n *Notifier) Send(alerts []alert.Alert) error {
	sum := summary.Build(alerts)
	payload := Payload{
		Timestamp:    time.Now().UTC(),
		OverallLevel: string(summary.OverallLevel(sum)),
		Total:        sum.Total,
		Critical:     sum.Critical,
		Warn:         sum.Warn,
		Ecosystems:   sum.Ecosystems,
		Alerts:       alerts,
	}

	switch n.cfg.Channel {
	case ChannelWebhook:
		return n.sendWebhook(payload)
	default:
		return n.sendStdout(payload)
	}
}

func (n *Notifier) sendStdout(p Payload) error {
	fmt.Fprintf(n.out, "[depwatch] %s | total=%d critical=%d warn=%d\n",
		p.OverallLevel, p.Total, p.Critical, p.Warn)
	return nil
}

func (n *Notifier) sendWebhook(p Payload) error {
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}
	resp, err := n.client.Post(n.cfg.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
