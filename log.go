package main

import "go.uber.org/zap"

func initLog() {
	zapcfg := zap.NewProductionConfig()
	zapcfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	zapcfg.DisableStacktrace = true
	zapLogger, _ := zapcfg.Build()
	zap.ReplaceGlobals(zapLogger)
}
