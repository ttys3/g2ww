package main

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sync"
)

// https://grafana.com/docs/grafana/latest/alerting/old-alerting/notifications/#webhook

type Hook struct {
	DashboardID int `json:"dashboardId"`
	Evalmatches []struct {
		Value  int    `json:"value"`
		Metric string `json:"metric"`
		Tags   struct {
		} `json:"tags"`
	} `json:"evalMatches"`
	ImageURL string `json:"imageUrl"`
	Message  string `json:"message"`
	OrgID    int    `json:"orgId"`
	PanelID  int    `json:"panelId"`
	RuleID   int    `json:"ruleId"`
	RuleName string `json:"ruleName"`
	RuleURL  string `json:"ruleUrl"`
	State    string `json:"state"`
	Tags     struct {
		TagName string `json:"tag name"`
	} `json:"tags"`
	Title string `json:"title"`
}

var sentCountLock sync.RWMutex
var sentCount int = 0

func GwStat(c echo.Context) error {
	sentCountLock.RLock()
	sc := sentCount
	sentCountLock.RUnlock()
	statMsg := fmt.Sprintf("G2WW Server is running! \nParsed & forwarded %d messages to WeChat Work!", sc)
	return c.String(http.StatusOK, statMsg)
}

func GwWorker(c echo.Context) error {
	var h Hook
	if err := c.Bind(&h); err != nil {
		zap.S().Errorw("error parsing json", "err", err)
		return fmt.Errorf("error parsing json, err=%v", err)
	}

	// Send to WeChat Work

	// {
	// 	"msgtype": "news",
	// 	"news": {
	// 	  "articles": [
	// 		{
	// 		  "title": "%s",
	// 		  "description": "%s",
	// 		  "url": "%s",
	// 		  "picurl": "%s"
	// 		}
	// 	  ]
	// 	}
	//   }
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + c.Param("key")

	msgStr := fmt.Sprintf(`
		{
			"msgtype": "news",
			"news": {
			  "articles": [
				{
				  "title": "%s",
				  "description": "%s",
				  "url": "%s",
				  "picurl": "%s"
				}
			  ]
			}
		  }
		`, h.Title, h.Message, h.RuleURL, h.ImageURL)
	zap.S().Infow("begin send data to wechat work webhook", "json", msgStr)
	jsonStr := []byte(msgStr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error sending to WeChat Work API")
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	sentCountLock.Lock()
	sentCount++
	sentCountLock.Unlock()
	zap.S().Infow("send data to wechat work webhook success")
	return c.String(http.StatusOK, string(rspBody))
}
