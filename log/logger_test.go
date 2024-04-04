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
	"regexp"
	"strconv"
	"testing"
)

func testOptions() *Options {
	return DefaultOptions()
}

func TestBasic(t *testing.T) {
	l := &logger{
		name:       "beep",
		callerSkip: 0,
	}
	l.SetOutputLevel(InfoLevel)
	l.SetStackTraceLevel(NoneLevel)
	l.SetLogCallers(false)

	cases := []struct {
		f          func()
		pat        string
		json       bool
		caller     bool
		wantExit   bool
		stackLevel Level
	}{
		{
			f:   func() { l.Debug("Hello") },
			pat: timePattern + "\tdebug\tbeep\tHello",
		},
		{
			f:   func() { l.Debugf("Hello") },
			pat: timePattern + "\tdebug\tbeep\tHello",
		},
		{
			f:   func() { l.Debugf("%s", "Hello") },
			pat: timePattern + "\tdebug\tbeep\tHello",
		},
		{
			f:      func() { l.Debug("Hello") },
			pat:    timePattern + "\tdebug\tbeep\tlog/logger_test.go:.*\tHello",
			caller: true,
		},
		{
			f: func() { l.Debug("Hello") },
			pat: "{\"level\":\"debug\",\"time\":\"" + timePattern + "\",\"scope\":\"beep\",\"caller\":\"log/logger_test.go:.*\",\"msg\":\"Hello\"," +
				"\"stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},

		{
			f:   func() { l.Info("Hello") },
			pat: timePattern + "\tinfo\tbeep\tHello",
		},
		{
			f:   func() { l.Infof("Hello") },
			pat: timePattern + "\tinfo\tbeep\tHello",
		},
		{
			f:   func() { l.Infof("%s", "Hello") },
			pat: timePattern + "\tinfo\tbeep\tHello",
		},
		{
			f: func() { l.Info("Hello") },
			pat: "{\"level\":\"info\",\"time\":\"" + timePattern + "\",\"scope\":\"beep\",\"caller\":\"log/logger_test.go:.*\",\"msg\":\"Hello\"," +
				"\"stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},

		{
			f:   func() { l.Warn("Hello") },
			pat: timePattern + "\twarn\tbeep\tHello",
		},
		{
			f:   func() { l.Warnf("Hello") },
			pat: timePattern + "\twarn\tbeep\tHello",
		},
		{
			f:   func() { l.Warnf("%s", "Hello") },
			pat: timePattern + "\twarn\tbeep\tHello",
		},
		{
			f: func() { l.Warn("Hello") },
			pat: "{\"level\":\"warn\",\"time\":\"" + timePattern + "\",\"scope\":\"beep\",\"caller\":\"log/logger_test.go:.*\",\"msg\":\"Hello\"," +
				"\"stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},

		{
			f:   func() { l.Error("Hello") },
			pat: timePattern + "\terror\tbeep\tHello",
		},
		{
			f:   func() { l.Errorf("Hello") },
			pat: timePattern + "\terror\tbeep\tHello",
		},
		{
			f:   func() { l.Errorf("%s", "Hello") },
			pat: timePattern + "\terror\tbeep\tHello",
		},
		{
			f: func() { l.Error("Hello") },
			pat: "{\"level\":\"error\",\"time\":\"" + timePattern + "\",\"scope\":\"beep\",\"caller\":\"log/logger_test.go:.*\"," +
				"\"msg\":\"Hello\"," +
				"\"stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},

		{
			f:        func() { l.Fatal("Hello") },
			pat:      timePattern + "\tfatal\tbeep\tHello",
			wantExit: true,
		},
		{
			f:        func() { l.Fatalf("Hello") },
			pat:      timePattern + "\tfatal\tbeep\tHello",
			wantExit: true,
		},
		{
			f:        func() { l.Fatalf("%s", "Hello") },
			pat:      timePattern + "\tfatal\tbeep\tHello",
			wantExit: true,
		},
		{
			f: func() { l.Fatal("Hello") },
			pat: "{\"level\":\"fatal\",\"time\":\"" + timePattern + "\",\"scope\":\"beep\",\"caller\":\"log/logger_test.go:.*\"," +
				"\"msg\":\"Hello\"," +
				"\"stack\":\".*\"}",
			json:       true,
			caller:     true,
			wantExit:   true,
			stackLevel: DebugLevel,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var exitCalled bool
			lines, err := captureStdout(func() {
				o := testOptions()
				o.JSONEncoding = c.json

				if err := Configure(o); err != nil {
					t.Errorf("Got err '%v', expecting success", err)
				}

				pt := funcs.Load().(functionTable)
				pt.exitProcess = func(_ int) {
					exitCalled = true
				}
				funcs.Store(pt)

				l.SetOutputLevel(DebugLevel)
				l.SetStackTraceLevel(c.stackLevel)
				l.SetLogCallers(c.caller)

				c.f()
				_ = Sync()
			})

			if exitCalled != c.wantExit {
				var verb string
				if c.wantExit {
					verb = " never"
				}
				t.Errorf("ol.Exit%s called", verb)
			}

			if err != nil {
				t.Errorf("Got error '%v', expected success", err)
			}

			if match, _ := regexp.MatchString(c.pat, lines[0]); !match {
				t.Errorf("Got '%v',\nexpected a match with '%v'", lines[0], c.pat)
			}
		})
	}
}
