package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func setup() {
	logrus.SetOutput(os.Stdout)

	ll := logrus.WarnLevel

	if lv, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch lv {
		case "TRACE":
			ll = logrus.TraceLevel
		case "DEBUG":
			ll = logrus.DebugLevel
		case "INFO":
			ll = logrus.InfoLevel
		case "WARN":
			ll = logrus.WarnLevel
		case "PANIC":
			ll = logrus.PanicLevel
		case "ERROR":
			ll = logrus.ErrorLevel
		case "FATAL":
			ll = logrus.FatalLevel
		}
	}

	logrus.SetLevel(ll)
}
