package main

import (
	"time"
)

func getSuffix(format string) string {
	currentTime := time.Now()
	return currentTime.Format(format)
}
