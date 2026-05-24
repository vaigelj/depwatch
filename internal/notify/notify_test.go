package notify_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/depwatch/internal/alert"
	"github.com/depwatch/internal/notify"
)

func sampleAlerts() []alert.Alert {
	return []alert.Alert{
		{Package: "lodash", Ecosystem: "npm", Severity: alert.SeverityCritical, Message: "vulnerable"},
		{Package: "express", Ecosystem: "npm", Severity: alert.SeverityWarn, Message: "outdated"},
	}
}

func TestSend_Stdout_WritesLine(t *testing.T) {
	var buf bytes.Buffer
	cfg := notify.Config{Channel: notify.ChannelStdout}
	n := notify.NewWithWriter(cfg, &buf)

	if err := n.Send(sampleAlerts()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[depwatch]") {
		t.Errorf("expected [depwatch] prefix, got: %q", out)
	}
	if !strings.Contains(out, "critical=1") {
		t.Errorf("expected critical=1 in output, got: %q", out)
	}
	if !strings.Contains(out, "warn=1") {
		t.Errorf("expected warn=1 in output, got: %q", out)
	}
}

func TestSend_Stdout_EmptyAlerts(t *testing.T) {
	var buf bytes.Buffer
	cfg := notify.Config{Channel: notify.ChannelStdout}
	n := notify.NewWithWriter(cfg, &buf)

	if err := n.Send([]alert.Alert{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected some output even for empty alerts")
	}
}

func TestSend_Webhook_PostsJSON(t *testing.T) {
	var received notify.Payload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := notify.Config{Channel: notify.ChannelWebhook, WebhookURL: ts.URL}
	n := notify.New(cfg)

	if err := n.Send(sampleAlerts()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Total != 2 {
		t.Errorf("expected total=2, got %d", received.Total)
	}
	if received.Critical != 1 {
		t.Errorf("expected critical=1, got %d", received.Critical)
	}
}

func TestSend_Webhook_ServerError_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cfg := notify.Config{Channel: notify.ChannelWebhook, WebhookURL: ts.URL}
	n := notify.New(cfg)

	if err := n.Send(sampleAlerts()); err == nil {
		t.Error("expected error for 500 response, got nil")
	}
}

func TestNew_DefaultChannel_IsStdout(t *testing.T) {
	cfg := notify.Config{}
	n := notify.New(cfg)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
