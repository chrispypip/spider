package spider

import (
	"errors"
	"fmt"
	"io"
	"log/syslog"
	"os"

	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/sirupsen/logrus/hooks/writer"
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
	LogTTYFormatter = iota
	LogTextFormatter
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

func SetLogOutput(out io.Writer) {
	log.SetOutput(out)
}

func AddLogFile(path string, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %s", path, err)
	}
	log.AddHook(&writer.Hook{
		Writer: file,
		LogLevels: []log.Level{
			log.TraceLevel,
			log.DebugLevel,
			log.InfoLevel,
			log.WarnLevel,
			log.ErrorLevel,
			log.FatalLevel,
			log.PanicLevel,
		},
	})
	return file, nil
}

func RemoveLogFile(file *os.File) {
	_ = file.Close()
}

func AddLogToSyslog(network, raddr string, priority syslog.Priority, tag string) error {
	hook, err := lSyslog.NewSyslogHook(network, raddr, priority, tag)
	if err != nil {
		return fmt.Errorf("failed to open syslog: %s", err)
	}
	log.AddHook(hook)
	return nil
}

func SetLogFormatter(logFormatter LogFormatter, enableColors, enableTimestamp, enablePrettyPrint bool) error {
	switch logFormatter {
	case LogTTYFormatter:
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
		return nil
	case LogTextFormatter:
		log.SetFormatter(&log.TextFormatter{
			DisableColors:    !enableColors,
			DisableTimestamp: !enableTimestamp,
		})
		return nil
	case LogJSONFormatter:
		log.SetFormatter(&log.JSONFormatter{
			DisableTimestamp: !enableTimestamp,
			PrettyPrint:      enablePrettyPrint,
		})
		return nil
	default:
		return errors.New("invalid LogFormatter specified; valid values are: LogTextFormatter, LogTTYFormatter, LogJSONFormatter")
	}
}
