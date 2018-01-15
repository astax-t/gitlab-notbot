package main

import (
	golog "log"
)

const (
	LOG_FATAL   = 0
	LOG_MESSAGE = 1
	LOG_ERROR   = 1
	LOG_INFO    = 2
	LOG_DEBUG   = 3
)

func log(level int, message string, err error)  {
	if level > config.LogLevel {
		return
	}

	if level == LOG_FATAL {
		golog.Fatal(message, err)
	}

	if err != nil {
		golog.Println(message, err)
	} else {
		golog.Println(message)
	}
}
