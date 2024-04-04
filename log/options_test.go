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
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestOptions(t *testing.T) {
	cases := []struct {
		cmdLine string
		result  Options
	}{
		{"--log_as_json", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			JSONEncoding:       true,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "info",
			StackTraceLevel:    "none",
			LogCaller:          false,
		}},

		{"--log_path stdout", Options{
			OutputPath:         "stdout",
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "info",
			StackTraceLevel:    "none",
			LogCaller:          false,
		}},

		{"--log_caller", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "info",
			StackTraceLevel:    "none",
			LogCaller:          true,
		}},

		{"--log_stacktrace_level debug", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "info",
			StackTraceLevel:    "debug",
			LogCaller:          false,
		}},

		{"--log_stacktrace_level info", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "info",
			StackTraceLevel:    "info",
			LogCaller:          false,
		}},

		{"--log_stacktrace_level warn", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "info",
			StackTraceLevel:    "warn",
			LogCaller:          false,
		}},

		{"--log_output_level debug", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "debug",
			StackTraceLevel:    "none",
			LogCaller:          false,
		}},

		{"--log_output_level warn", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "warn",
			StackTraceLevel:    "none",
			LogCaller:          false,
		}},

		{"--log_rotate_path foobar", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotateOutputPath:   "foobar",
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "info",
			StackTraceLevel:    "none",
			LogCaller:          false,
		}},

		{"--log_rotate_max_age 1234", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     1234,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "info",
			StackTraceLevel:    "none",
			LogCaller:          false,
		}},

		{"--log_rotate_max_size 1234", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    1234,
			RotationMaxBackups: DefaultRotationMaxBackups,
			OutputLevel:        "info",
			StackTraceLevel:    "none",
			LogCaller:          false,
		}},

		{"--log_rotate_max_backups 1234", Options{
			OutputPath:         DefaultOutputPath,
			ErrorOutputPath:    DefaultErrorOutputPath,
			RotationMaxAge:     DefaultRotationMaxAge,
			RotationMaxSize:    DefaultRotationMaxSize,
			RotationMaxBackups: 1234,
			OutputLevel:        "info",
			StackTraceLevel:    "none",
			LogCaller:          false,
		}},
	}

	for j := 0; j < 2; j++ {
		for i, c := range cases {
			t.Run(strconv.Itoa(j*100+i), func(t *testing.T) {
				o := DefaultOptions()
				cmd := &cobra.Command{}
				o.AddFlags(cmd.Flags())
				cmd.SetArgs(strings.Split(c.cmdLine, " "))

				if err := cmd.Execute(); err != nil {
					t.Errorf("Got %v, expecting success", err)
				}

				if !reflect.DeepEqual(c.result, *o) {
					t.Errorf("Got %v, expected %v", *o, c.result)
				}
			})
		}
	}
}
