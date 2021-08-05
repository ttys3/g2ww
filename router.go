package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

var sentCountLock sync.RWMutex
var sentCount int = 0

func GwStat(c echo.Context) error {
	sentCountLock.RLock()
	sc := sentCount
	sentCountLock.RUnlock()
	statMsg := fmt.Sprintf("G2WW Server is running! \nParsed & forwarded %d messages to WeChat Work!", sc)
	return c.String(http.StatusOK, statMsg)
}

func ProxyWebhook(c echo.Context) error {
	var h GrafanaAlertMsg
	if err := c.Bind(&h); err != nil {
		zap.S().Errorw("error parsing json", "err", err)
		return fmt.Errorf("error parsing json, err=%v", err)
	}

	key := c.Param("key")
	if key == "" {
		return c.String(http.StatusBadRequest, "key can not bey empty")
	}

	rspBody, err := SendWebhook(key, &h)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	sentCountLock.Lock()
	sentCount++
	sentCountLock.Unlock()

	zap.S().Infow("send data to wechat work webhook success")
	return c.JSON(http.StatusOK, rspBody)
}
