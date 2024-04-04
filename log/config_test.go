package log

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const timePattern = "[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9].[0-9][0-9][0-9][0-9][0-9][0-9]Z"

// Runs the given function while capturing everything sent to stdout
func captureStdout(f func()) ([]string, error) {
	tf, err := os.CreateTemp("", "log_test")
	if err != nil {
		return nil, err
	}

	old := os.Stdout
	os.Stdout = tf

	f()

	os.Stdout = old
	path := tf.Name()
	_ = tf.Sync()
	_ = tf.Close()

	content, err := os.ReadFile(path)
	_ = os.Remove(path)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}

func TestOverrides(t *testing.T) {
	o := DefaultOptions()
	o.OutputLevel = "debug"
	if err := Configure(o); err != nil {
		t.Errorf("Expecting success, got %v", err)
	} else if defaultLogger.GetOutputLevel() != DebugLevel {
		t.Errorf("Expecting DebugLevel, got %v", defaultLogger.GetOutputLevel())
	}

	o = DefaultOptions()
	o.StackTraceLevel = "debug"
	if err := Configure(o); err != nil {
		t.Errorf("Expecting success, got %v", err)
	} else if defaultLogger.GetStackTraceLevel() != DebugLevel {
		t.Errorf("Expecting DebugLevel, got %v", defaultLogger.GetStackTraceLevel())
	}

	o = DefaultOptions()
	o.LogCaller = true
	if err := Configure(o); err != nil {
		t.Errorf("Expecting success, got %v", err)
	} else if !defaultLogger.GetLogCallers() {
		t.Error("Expecting true, got false")
	}
}

func TestCapture(t *testing.T) {
	lines, _ := captureStdout(func() {
		o := DefaultOptions()
		o.LogCaller = true
		o.OutputLevel = "debug"
		_ = Configure(o)

		// output to the plain golang "log" package
		log.Println("golang")

		// output directly to zap
		zap.L().Error("zap-error")
		zap.L().Warn("zap-warn")
		zap.L().Info("zap-info")
		zap.L().Debug("zap-debug")

		l := zap.L().With(zap.String("a", "b"))
		l.Error("zap-with")

		entry := zapcore.Entry{
			Message: "zap-write",
			Level:   zapcore.ErrorLevel,
		}
		_ = zap.L().Core().Write(entry, nil)

		defaultLogger.SetOutputLevel(NoneLevel)

		// all these get thrown out since the level is set to none
		log.Println("golang-2")
		zap.L().Error("zap-error-2")
		zap.L().Warn("zap-warn-2")
		zap.L().Info("zap-info-2")
		zap.L().Debug("zap-debug-2")
	})

	patterns := []string{
		timePattern + "\tinfo\tlog/config_test.go:.*\tgolang",
		timePattern + "\terror\tlog/config_test.go:.*\tzap-error",
		timePattern + "\twarn\tlog/config_test.go:.*\tzap-warn",
		timePattern + "\tinfo\tlog/config_test.go:.*\tzap-info",
		timePattern + "\tdebug\tlog/config_test.go:.*\tzap-debug",
		timePattern + "\terror\tlog/config_test.go:.*\tzap-with",
		"error\tzap-write",
		"",
	}

	if len(lines) > len(patterns) {
		t.Errorf("Expecting %d lines of output, but got %d", len(patterns), len(lines))

		for i := len(patterns); i < len(lines); i++ {
			t.Errorf("  Extra line of output: %s", lines[i])
		}
	}

	for i, pat := range patterns {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			match, _ := regexp.MatchString(pat, lines[i])
			if !match {
				t.Errorf("Got '%s', expecting to match '%s'", lines[i], pat)
			}
		})
	}

	lines, _ = captureStdout(func() {
		o := DefaultOptions()
		o.StackTraceLevel = "debug"
		o.OutputLevel = "debug"
		_ = Configure(o)
		log.Println("golang")
	})

	for _, line := range lines {
		// see if the captured output contains the current file name
		if strings.Contains(line, "config_test.go") {
			return
		}
	}

	t.Error("Could not find stack trace info in output")
}
