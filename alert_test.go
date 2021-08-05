package main

import (
	"testing"
	"time"
)

func TestSendWebhook(t *testing.T) {
	key := "your key here"

	h := &GrafanaAlertMsg{
		Receiver:    "",
		Status:      "firing",
		Externalurl:     "https://bing.com",
		Version:         "",
		Groupkey:        "",
		Title:           "测试alert",
		State:           "alerting",
		CommonAnnotations: CommonAnnotations{
			"too many invalid client",
		},
		Message:         "**Firing**\n\nLabels:\n - alertname = invalid_client\n - to = group\nAnnotations:\n - description = too many invalid client\nSource: ....",
		Alerts: []GrafanaAlertItem{
			GrafanaAlertItem{
				Status: "firing",
				Annotations: struct {
    Description string `json:"description"`
}{"this is Annotations description"},
				StartsAt:     time.Now(),
				EndsAt:       time.Now().Add(time.Hour),
				GeneratorURL: "",
				Fingerprint:  "",
				SilenceURL:   "",
				DashboardURL: "",
				PanelURL:     "https://bing.com",
				ValueString:  "32",
			},
		},
	}

	rsp, err := SendWebhook(key, h)
	t.Logf("err=%v rsp=%+v", err, rsp)
}