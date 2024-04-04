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

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	DefaultLoggerName         = "default"
	DefaultOutputLevel        = InfoLevel
	DefaultStackTraceLevel    = NoneLevel
	DefaultOutputPath         = "stdout"
	DefaultErrorOutputPath    = "stderr"
	DefaultRotationMaxAge     = 30
	DefaultRotationMaxSize    = 100 * 1024 * 1024
	DefaultRotationMaxBackups = 1000
)

// Level is an enumeration of all supported log levels.
type Level int32

const (
	// NoneLevel disables logging
	NoneLevel Level = iota
	// FatalLevel enables fatal level logging
	FatalLevel
	// ErrorLevel enables error level logging
	ErrorLevel
	// WarnLevel enables warn level logging
	WarnLevel
	// InfoLevel enables info level logging
	InfoLevel
	// DebugLevel enables debug level logging
	DebugLevel
)

var levelToString = map[Level]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warn",
	ErrorLevel: "error",
	FatalLevel: "fatal",
	NoneLevel:  "none",
}

var stringToLevel = map[string]Level{
	"debug": DebugLevel,
	"info":  InfoLevel,
	"warn":  WarnLevel,
	"error": ErrorLevel,
	"fatal": FatalLevel,
	"none":  NoneLevel,
}

var toLevel = map[zapcore.Level]Level{
	zapcore.FatalLevel: FatalLevel,
	zapcore.ErrorLevel: ErrorLevel,
	zapcore.WarnLevel:  WarnLevel,
	zapcore.InfoLevel:  InfoLevel,
	zapcore.DebugLevel: DebugLevel,
}

var levelToZap = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	FatalLevel: zapcore.FatalLevel,
}

var levelListString = []string{"debug", "info", "warn", "error", "fatal", "none"}

// Options defines the set of options supported by component-base logging package.
type Options struct {
	// OutputPath is a file system path to write the log data to.
	// The special values stdout and stderr can be used to output to the
	// standard I/O streams. This defaults to stdout.
	OutputPath string

	// ErrorOutputPath is a file system path to write logger errors to.
	// The special values stdout and stderr can be used to output to the
	// standard I/O streams. This defaults to stderr.
	ErrorOutputPath string

	// RotateOutputPath is the path to a rotating log file. This file should
	// be automatically rotated over time, based on the rotation parameters such
	// as RotationMaxSize and RotationMaxAge. The default is to not rotate.
	//
	// This path is used as a foundational path. This is where log output is normally
	// saved. When a rotation needs to take place because the file got too big or too
	// old, then the file is renamed by appending a timestamp to the name. Such renamed
	// files are called backups. Once a backup has been created,
	// output resumes to this path.
	RotateOutputPath string

	// RotationMaxSize is the maximum size in megabytes of a log file before it gets
	// rotated. It defaults to 100 megabytes.
	RotationMaxSize int

	// RotationMaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename. Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is to remove log files
	// older than 30 days.
	RotationMaxAge int

	// RotationMaxBackups is the maximum number of old log files to retain.  The default
	// is to retain at most 1000 logs.
	RotationMaxBackups int

	// JSONEncoding controls whether the log is formatted as JSON.
	JSONEncoding bool

	// OutputLevel controls the log level.
	OutputLevel string

	// StackTraceLevel controls the log level for stack trace.
	StackTraceLevel string

	// LogCaller controls whether to log the caller of a logging function
	LogCaller bool
}

// DefaultOptions returns a new set of options, initialized to the defaults
func DefaultOptions() *Options {
	return &Options{
		OutputPath:         DefaultOutputPath,
		ErrorOutputPath:    DefaultErrorOutputPath,
		RotationMaxSize:    DefaultRotationMaxSize,
		RotationMaxAge:     DefaultRotationMaxAge,
		RotationMaxBackups: DefaultRotationMaxBackups,
		OutputLevel:        levelToString[InfoLevel],
		StackTraceLevel:    levelToString[NoneLevel],
		LogCaller:          false,
	}
}

// AddFlags add logging-format flag.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.OutputPath, "log_path", o.OutputPath,
		"The file path where to output the log. This can be any path as well as the special values stdout and stderr")

	fs.StringVar(&o.RotateOutputPath, "log_rotate_path", o.RotateOutputPath,
		"The file path for the optional rotating log file")

	fs.IntVar(&o.RotationMaxAge, "log_rotate_max_age", o.RotationMaxAge,
		"The maximum age in days of a log file beyond which the file is rotated (0 indicates no limit)")

	fs.IntVar(&o.RotationMaxSize, "log_rotate_max_size", o.RotationMaxSize,
		"The maximum size in megabytes of a log file beyond which the file is rotated")

	fs.IntVar(&o.RotationMaxBackups, "log_rotate_max_backups", o.RotationMaxBackups,
		"The maximum number of log file backups to keep before older files are deleted (0 indicates no limit)")

	fs.BoolVar(&o.JSONEncoding, "log_as_json", o.JSONEncoding,
		"Whether to format output as JSON or in plain console-friendly format")

	fs.StringVar(&o.OutputLevel, "log_output_level", o.OutputLevel,
		fmt.Sprintf("The minimum logging level of messages to output,  can be one of %s",
			levelListString))

	fs.StringVar(&o.StackTraceLevel, "log_stacktrace_level", o.StackTraceLevel,
		fmt.Sprintf("The minimum logging level at which stack traces are captured, can be one of %s",
			levelListString))

	fs.BoolVar(&o.LogCaller, "log_caller", o.LogCaller, "Whether to log the caller of a logging function or not")
}
