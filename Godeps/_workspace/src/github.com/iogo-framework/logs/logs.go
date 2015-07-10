// Package logs provide basic and unified logging functions.
package logs

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

// DefaultLogger provide a simple interface et sane defaults.
// By default, it will not log debug infos.
var DefaultLogger = &Logger{Level: InfoLevel}

const (
	CriticalLevel = iota
	ErrorLevel
	WarningLevel
	NoticeLevel
	InfoLevel
	DebugLevel
)

type LevelInfo struct {
	Str   string
	Color color.Attribute
}

var Levels = map[int]LevelInfo{
	CriticalLevel: {Str: "CRIT", Color: color.FgMagenta},
	ErrorLevel:    {Str: "ERRO", Color: color.FgRed},
	WarningLevel:  {Str: "WARN", Color: color.FgYellow},
	NoticeLevel:   {Str: "NOTI", Color: color.FgGreen},
	InfoLevel:     {Str: "INFO", Color: color.Reset},
	DebugLevel:    {Str: "DEBU", Color: color.FgCyan},
}

// Logger limit the output logs to a specified level.
type Logger struct {
	Level int
}

func logs(level int, maxLevel int, format interface{}, v ...interface{}) {
	if level > maxLevel {
		return
	}
	levelStr := Levels[level].Str
	colorPrint := color.New(Levels[level].Color).SprintFunc()
	switch format.(type) {
	case string:
		log.Printf("%s %s\n", colorPrint(levelStr), fmt.Sprintf(format.(string), v...))
	case error:
		log.Printf("%s %s\n", colorPrint(levelStr), fmt.Sprint(format.(error)))
	default:
		log.Printf("%s %v\n", colorPrint(levelStr), format)
	}
}

func (l Logger) Debug(format interface{}, v ...interface{}) {
	logs(DebugLevel, l.Level, format, v...)
}

func (l Logger) Info(format interface{}, v ...interface{}) {
	logs(InfoLevel, l.Level, format, v...)
}

func (l Logger) Notice(format interface{}, v ...interface{}) {
	logs(NoticeLevel, l.Level, format, v...)
}

func (l Logger) Warning(format interface{}, v ...interface{}) {
	logs(WarningLevel, l.Level, format, v...)
}

func (l Logger) Error(format interface{}, v ...interface{}) {
	logs(ErrorLevel, l.Level, format, v...)
}

func (l Logger) Critical(format interface{}, v ...interface{}) {
	logs(CriticalLevel, l.Level, format, v...)
}

// Level set the DefaultLogger log level
func Level(level int) {
	DefaultLogger.Level = level
}

// Debug logger.
var Debug = func(format interface{}, v ...interface{}) { DefaultLogger.Debug(format, v...) }

// Info logger.
var Info = func(format interface{}, v ...interface{}) { DefaultLogger.Info(format, v...) }

// Notice logger.
var Notice = func(format interface{}, v ...interface{}) { DefaultLogger.Notice(format, v...) }

// Warning logger.
var Warning = func(format interface{}, v ...interface{}) { DefaultLogger.Warning(format, v...) }

// Error logger.
var Error = func(format interface{}, v ...interface{}) { DefaultLogger.Error(format, v...) }

// Critical logger.
var Critical = func(format interface{}, v ...interface{}) { DefaultLogger.Critical(format, v...) }
