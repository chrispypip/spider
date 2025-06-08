package spider

import (
	log "github.com/sirupsen/logrus"
)

type LogLevel uint8

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelFatal
	LogLevelPanic
)

type LogFormatter uint8

const (
	LogTextFormatter = iota
	LogTTYFormatter
	LogJSONFormatter
)

func SetLogLevel(level LogLevel) {
	switch level {
	case LogLevelTrace:
		log.SetLevel(log.TraceLevel)
	case LogLevelDebug:
		log.SetLevel(log.DebugLevel)
	case LogLevelInfo:
		log.SetLevel(log.InfoLevel)
	case LogLevelWarning:
		log.SetLevel(log.WarnLevel)
	case LogLevelError:
		log.SetLevel(log.ErrorLevel)
	case LogLevelFatal:
		log.SetLevel(log.FatalLevel)
	case LogLevelPanic:
		log.SetLevel(log.PanicLevel)
	}
}

func GetLogLevel() LogLevel {
	level := log.GetLevel()
	switch level {
	case log.TraceLevel:
		return LogLevelTrace
	case log.DebugLevel:
		return LogLevelDebug
	case log.InfoLevel:
		return LogLevelInfo
	case log.WarnLevel:
		return LogLevelWarning
	case log.ErrorLevel:
		return LogLevelError
	case log.FatalLevel:
		return LogLevelFatal
	case log.PanicLevel:
		return LogLevelPanic
	}
	return LogLevelInfo
}

func SetLogFormatter(logFormatter LogFormatter) {
	switch logFormatter {
	case LogTextFormatter:
		log.SetFormatter(&log.TextFormatter{})
	case LogTTYFormatter:
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	case LogJSONFormatter:
		log.SetFormatter(&log.JSONFormatter{})
	default:
		panic("Invalid LogFormatter specified")
	}
}
