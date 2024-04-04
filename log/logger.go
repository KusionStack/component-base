// Copyright 2024 KusionStack Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// callerSkipOffset is how many callers to pop off the stack to determine the caller function locality, used for
// adding file/line number to log output.
const callerSkipOffset = 3

var defaultLogger *logger

// logger collects all the global state of the logging setup.
type logger struct {
	name       string
	callerSkip int

	outputLevel     atomic.Value
	stackTraceLevel atomic.Value
	logCallers      atomic.Value
}

// Info outputs a message at info level.
func (l *logger) Info(field any) {
	if l.GetOutputLevel() >= InfoLevel {
		l.output(zapcore.InfoLevel, fmt.Sprint(field))
	}
}

// Infof uses fmt.Sprintf to construct and outputs a message at info level.
func (l *logger) Infof(format string, fields ...any) {
	if l.GetOutputLevel() >= InfoLevel {
		msg := maybeSprintf(format, fields...)
		l.output(zapcore.InfoLevel, msg)
	}
}

// InfoEnabled returns whether output of messages using this logger is currently enabled for info-level output.
func (l *logger) InfoEnabled() bool {
	return l.GetOutputLevel() >= InfoLevel
}

// Debug outputs a message at debug level.
func (l *logger) Debug(field any) {
	if l.GetOutputLevel() >= DebugLevel {
		l.output(zapcore.DebugLevel, fmt.Sprint(field))
	}
}

// Debugf uses fmt.Sprintf to construct and outputs a message at debug level.
func (l *logger) Debugf(format string, fields ...any) {
	if l.GetOutputLevel() >= DebugLevel {
		msg := maybeSprintf(format, fields...)
		l.output(zapcore.DebugLevel, msg)
	}
}

// DebugEnabled returns whether output of messages using this logger is currently enabled for debug-level output.
func (l *logger) DebugEnabled() bool {
	return l.GetOutputLevel() >= DebugLevel
}

// Warn outputs a message at warn level.
func (l *logger) Warn(field any) {
	if l.GetOutputLevel() >= WarnLevel {
		l.output(zapcore.WarnLevel, fmt.Sprint(field))
	}
}

// Warnf uses fmt.Sprintf to construct and outputs a message at warn level.
func (l *logger) Warnf(format string, fields ...any) {
	if l.GetOutputLevel() >= WarnLevel {
		msg := maybeSprintf(format, fields...)
		l.output(zapcore.WarnLevel, msg)
	}
}

// WarnEnabled returns whether output of messages using this logger is currently enabled for warn-level output.
func (l *logger) WarnEnabled() bool {
	return l.GetOutputLevel() >= WarnLevel
}

// Error outputs a message at error level.
func (l *logger) Error(field any) {
	if l.GetOutputLevel() >= ErrorLevel {
		l.output(zapcore.ErrorLevel, fmt.Sprint(field))
	}
}

// Errorf uses fmt.Sprintf to construct and outputs a message at error level.
func (l *logger) Errorf(format string, fields ...any) {
	if l.GetOutputLevel() >= ErrorLevel {
		msg := maybeSprintf(format, fields...)
		l.output(zapcore.ErrorLevel, msg)
	}
}

// ErrorEnabled returns whether output of messages using this logger is currently enabled for error-level output.
func (l *logger) ErrorEnabled() bool {
	return l.GetOutputLevel() >= ErrorLevel
}

// Fatal outputs a message at fatal level.
func (l *logger) Fatal(field any) {
	if l.GetOutputLevel() >= FatalLevel {
		l.output(zapcore.FatalLevel, fmt.Sprint(field))
	}
}

// Fatalf uses fmt.Sprintf to construct and outputs a message at fatal level.
func (l *logger) Fatalf(format string, fields ...any) {
	if l.GetOutputLevel() >= FatalLevel {
		msg := maybeSprintf(format, fields...)
		l.output(zapcore.FatalLevel, msg)
	}
}

// FatalEnabled returns whether output of messages using this logger is currently enabled for fatal-level output.
func (l *logger) FatalEnabled() bool {
	return l.GetOutputLevel() >= FatalLevel
}

// SetOutputLevel adjusts the output level associated with this logger.
func (l *logger) SetOutputLevel(level Level) {
	l.outputLevel.Store(level)
}

// GetOutputLevel returns the output level associated with this logger.
func (l *logger) GetOutputLevel() Level {
	return l.outputLevel.Load().(Level)
}

// SetStackTraceLevel adjusts the stack tracing level associated with this logger.
func (l *logger) SetStackTraceLevel(level Level) {
	l.stackTraceLevel.Store(level)
}

// GetStackTraceLevel returns the stack tracing level associated with this logger.
func (l *logger) GetStackTraceLevel() Level {
	return l.stackTraceLevel.Load().(Level)
}

// SetLogCallers adjusts the output level associated with this logger.
func (l *logger) SetLogCallers(logCallers bool) {
	l.logCallers.Store(logCallers)
}

// GetLogCallers returns the output level associated with this logger.
func (l *logger) GetLogCallers() bool {
	return l.logCallers.Load().(bool)
}

// output writes the data to the log files.
func (l *logger) output(level zapcore.Level, msg string) {
	e := zapcore.Entry{
		Message: msg,
		Level:   level,
		Time:    time.Now(),
	}
	if l.name != DefaultLoggerName {
		e.LoggerName = l.name
	}

	if l.GetLogCallers() {
		e.Caller = zapcore.NewEntryCaller(runtime.Caller(l.callerSkip + callerSkipOffset))
	}

	thresh := toLevel[level]
	if l.GetStackTraceLevel() >= thresh {
		e.Stack = zap.Stack("").String
	}

	ft := funcs.Load().(functionTable)
	if ft.write != nil {
		if err := ft.write(e, nil); err != nil {
			_, _ = fmt.Fprintf(ft.errorSink, "%v log write error: %v\n", time.Now(), err)
			_ = ft.errorSink.Sync()
		}
	}
}

// Info logs to the INFO log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Info(field any) {
	defaultLogger.Info(field)
}

// Infof logs to the INFO log.
func Infof(format string, fields ...any) {
	defaultLogger.Infof(format, fields...)
}

// Debug logs to the DEBUG log.
func Debug(field any) {
	defaultLogger.Debug(field)
}

// Debugf logs to the DEBUG log.
func Debugf(format string, fields ...any) {
	defaultLogger.Debugf(format, fields...)
}

// Warn logs to the WARN log.
func Warn(field any) {
	defaultLogger.Warn(field)
}

// Warnf logs to the WARN log.
func Warnf(format string, fields ...any) {
	defaultLogger.Warnf(format, fields...)
}

// Error logs to the ERROR log.
func Error(field any) {
	defaultLogger.Error(field)
}

// Errorf logs to the ERROR log.
func Errorf(format string, fields ...any) {
	defaultLogger.Errorf(format, fields...)
}

// Fatal logs to the FATAL log.
func Fatal(field any) {
	defaultLogger.Fatal(field)
}

// Fatalf logs to the FATAL log.
func Fatalf(format string, fields ...any) {
	defaultLogger.Fatalf(format, fields...)
}

func maybeSprintf(format string, args ...any) string {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	return msg
}
