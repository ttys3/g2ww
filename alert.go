package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type GrafanaAlertItem struct {
	Status string `json:"status"`
	Labels struct {
		Alertname string `json:"alertname"`
		To        string `json:"to"`
	} `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
	} `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Fingerprint  string    `json:"fingerprint"`
	SilenceURL   string    `json:"silenceURL"`
	DashboardURL string    `json:"dashboardURL"`
	PanelURL     string    `json:"panelURL"`
	ValueString  string    `json:"valueString"`
}

type CommonAnnotations struct {
	Description string `json:"description"`
}

type GrafanaAlertMsg struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []GrafanaAlertItem`json:"alerts"`
	Grouplabels map[string]string `json:"groupLabels"`
	Commonlabels map[string]string `json:"commonLabels"`
	CommonAnnotations CommonAnnotations `json:"commonAnnotations"`
	Externalurl     string `json:"externalURL"`
	Version         string `json:"version"`
	Groupkey        string `json:"groupKey"`
	Truncatedalerts int    `json:"truncatedAlerts"`
	Title           string `json:"title"`
	State           string `json:"state"`
	Message         string `json:"message"`
}

// {"errcode":0,"errmsg":"ok"}
// {"errcode":40039,"errmsg":"Warning: wrong json format. invalid url size, hint: [xxxxxxx], from ip: xxxxxxx, more info at https://open.work.weixin.qq.com/devtool/query?e=40039"}
// url长度限制1024个字节 https://open.work.weixin.qq.com/devtool/query?e=40039

// WechatWorkWebhookRsp 响应
type WechatWorkWebhookRsp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type MsgContent struct {
	Content string `json:"content"`
}

type MarkdownMsg struct {
	Msgtype  string `json:"msgtype"`
	Markdown MsgContent `json:"markdown"`
}

const MsgTemplate = `
Status: %s
Title: %s
Description: %s
PanelURL: [点此查看面板图](%s)
State: %s
`

func SendWebhook(key string, h *GrafanaAlertMsg) (*WechatWorkWebhookRsp, error) {
	if len(h.Alerts) == 0 {
		return nil, fmt.Errorf("err GrafanaAlertMsg.Alerts is empty")
	}
	// Send to WeChat Work
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + key
	msg := MarkdownMsg{
		Msgtype: "markdown",
		Markdown: MsgContent{
			Content: fmt.Sprintf(MsgTemplate, h.Status, h.Title, h.CommonAnnotations.Description, h.Alerts[0].PanelURL, h.State),
		},
	}
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		zap.S().Errorw("json.Marshal failed", "msg", msg, "err", err)
		return nil, err
	}

	zap.S().Infow("begin send data to wechat work webhook", "json", msg)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending to WeChat Work API")
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rsp WechatWorkWebhookRsp
	if err := json.Unmarshal(rspBody, &rsp); err != nil {
		return nil, fmt.Errorf("error unmarshal response body to json, err=%v", err)
	}
	return &rsp, nil
}
