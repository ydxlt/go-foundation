package log

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	mu     sync.Mutex
	level  Level
	logger *log.Logger
}

var (
	// 全局实例
	std  = New(INFO)
	once sync.Once
)

// New 初始化
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0), // 输出到 stdout
	}
}

func SetMode(debugMode bool) {
	if debugMode {
		SetLevel(DEBUG)
	} else {
		SetLevel(INFO)
	}
}

// SetLevel 设置全局日志级别
func SetLevel(level Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.level = level
}

// 内部输出
func (l *Logger) logf(level Level, format string, args ...any) {
	if level < l.level {
		return
	}
	levelStr := [...]string{"DEBUG", "INFO", "WARN", "ERROR"}[level]
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	l.logger.Printf("[%s] [%s] %s", timestamp, levelStr, msg)
}

func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := [...]string{"DEBUG", "INFO", "WARN", "ERROR"}[level]
	l.logger.Printf("[%s] [%s] %s", timestamp, levelStr, msg)
}

//
// ========== 日志方法 ==========
//

func Debugf(format string, args ...any) { std.logf(DEBUG, format, args...) }
func Infof(format string, args ...any)  { std.logf(INFO, format, args...) }
func Warnf(format string, args ...any)  { std.logf(WARN, format, args...) }
func Errorf(format string, args ...any) { std.logf(ERROR, format, args...) }
func Debug(msg string)                  { std.log(DEBUG, msg) }
func Info(msg string)                   { std.log(INFO, msg) }
func Warn(msg string)                   { std.log(WARN, msg) }
func Error(msg string)                  { std.log(ERROR, msg) }
