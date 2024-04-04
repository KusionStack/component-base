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
	"os"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// The default encoder config
var defaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "scope",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stack",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeTime:     formatDate,
}

// functionTable contains functions that can be replaced in a test setting
type functionTable struct {
	write       func(ent zapcore.Entry, fields []zapcore.Field) error
	sync        func() error
	exitProcess func(code int)
	errorSink   zapcore.WriteSyncer
	close       func() error
}

// functions that can be replaced by tests
var funcs = &atomic.Value{}

func init() {
	// use our defaults for starters so that logging works even before everything is fully configured
	_ = Configure(DefaultOptions())
}

// prepZap sets up the core Zap loggers
func prepZap(options *Options) (zapcore.Core, func() zapcore.Core, zapcore.WriteSyncer, error) {
	var enc zapcore.Encoder
	encCfg := defaultEncoderConfig

	if options.JSONEncoding {
		enc = zapcore.NewJSONEncoder(encCfg)
	} else {
		enc = zapcore.NewConsoleEncoder(encCfg)
	}

	var rotaterSink zapcore.WriteSyncer
	if options.RotateOutputPath != "" {
		rotaterSink = zapcore.AddSync(&lumberjack.Logger{
			Filename:   options.RotateOutputPath,
			MaxSize:    options.RotationMaxSize,
			MaxBackups: options.RotationMaxBackups,
			MaxAge:     options.RotationMaxAge,
		})
	}

	errSink, closeErrorSink, err := zap.Open(options.OutputPath)
	if err != nil {
		return nil, nil, nil, err
	}

	var outputSink zapcore.WriteSyncer
	if len(options.OutputPath) > 0 {
		outputSink, _, err = zap.Open(options.OutputPath)
		if err != nil {
			closeErrorSink()
			return nil, nil, nil, err
		}
	}

	var sink zapcore.WriteSyncer
	if rotaterSink != nil && outputSink != nil {
		sink = zapcore.NewMultiWriteSyncer(outputSink, rotaterSink)
	} else if rotaterSink != nil {
		sink = rotaterSink
	} else {
		sink = outputSink
	}

	alwaysOn := zapcore.NewCore(enc, sink, zap.NewAtomicLevelAt(zapcore.DebugLevel))
	conditionallyOn := func() zapcore.Core {
		enabler := func(lvl zapcore.Level) bool {
			switch lvl {
			case zapcore.ErrorLevel:
				return defaultLogger.ErrorEnabled()
			case zapcore.WarnLevel:
				return defaultLogger.WarnEnabled()
			case zapcore.InfoLevel:
				return defaultLogger.InfoEnabled()
			}
			return defaultLogger.DebugEnabled()
		}
		return zapcore.NewCore(enc, sink, zap.LevelEnablerFunc(enabler))
	}
	return alwaysOn, conditionallyOn, errSink, nil
}

// Configure initializes a functional logging subsystem.
//
// You typically call this once at process startup.
// Once this call returns, the logging system is ready to accept data.
func Configure(opts *Options) error {
	if err := updateLogger(opts); err != nil {
		return err
	}

	baseLogger, logBuilder, errSink, err := prepZap(opts)
	if err != nil {
		return err
	}

	// construct function table
	ft := functionTable{
		write: func(ent zapcore.Entry, fields []zapcore.Field) error {
			err := baseLogger.Write(ent, fields)
			if ent.Level == zapcore.FatalLevel {
				funcs.Load().(functionTable).exitProcess(1)
			}

			return err
		},
		sync:        baseLogger.Sync,
		exitProcess: os.Exit,
		errorSink:   errSink,
		close: func() error {
			// best-effort to sync
			_ = baseLogger.Sync()
			return nil
		},
	}
	funcs.Store(ft)

	zapOptions := []zap.Option{
		zap.ErrorOutput(errSink),
		zap.AddCallerSkip(1),
	}

	if defaultLogger.GetLogCallers() {
		zapOptions = append(zapOptions, zap.AddCaller())
	}

	l := defaultLogger.GetStackTraceLevel()
	if l != NoneLevel {
		zapOptions = append(zapOptions, zap.AddStacktrace(levelToZap[l]))
	}

	defaultZapLogger := zap.New(logBuilder(), zapOptions...)

	// capture global zap logging and force it through our logger
	_ = zap.ReplaceGlobals(defaultZapLogger)

	// capture standard golang "log" package output and force it through our logger
	_ = zap.RedirectStdLog(defaultZapLogger)

	return nil
}

func updateLogger(opts *Options) error {
	if defaultLogger == nil {
		defaultLogger = &logger{
			callerSkip: 1,
		}
	}

	level, ok := stringToLevel[opts.OutputLevel]
	if !ok {
		return fmt.Errorf("invalid output level '%s'", opts.OutputLevel)
	}
	defaultLogger.SetOutputLevel(level)

	if len(opts.StackTraceLevel) != 0 {
		stackTraceLevel, ok := stringToLevel[opts.StackTraceLevel]
		if !ok {
			return fmt.Errorf("invalid stack trace level '%s'", opts.OutputLevel)
		}
		defaultLogger.SetStackTraceLevel(stackTraceLevel)
	}

	defaultLogger.SetLogCallers(opts.LogCaller)

	return nil
}

func formatDate(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.UTC()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000

	buf := make([]byte, 27)

	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = 'T'
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	buf[23] = byte((micros/100)%10) + '0'
	buf[24] = byte((micros/10)%10) + '0'
	buf[25] = byte((micros)%10) + '0'
	buf[26] = 'Z'

	enc.AppendString(string(buf))
}

// Sync flushes any buffered log entries.
// Processes should normally take care to call Sync before exiting.
func Sync() error {
	return funcs.Load().(functionTable).sync()
}

// Close implements io.Closer.
func Close() error {
	return funcs.Load().(functionTable).close()
}
