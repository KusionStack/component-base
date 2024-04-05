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

package browser

import (
	"net/url"
	"os/exec"
	"runtime"

	"github.com/pkg/browser"
)

// nativeBrowser should implement the Browser interface.
var _ Browser = nativeBrowser{}

// NewNativeBrowser creates and returns a Browser that will attempt to interact
// with the browser-launching mechanisms of the operating system where the
// program is currently running.
func NewNativeBrowser() Browser {
	return nativeBrowser{}
}

type nativeBrowser struct{}

// OpenURL opens given url in a new browser tab.
func (b nativeBrowser) OpenURL(urlstr string) error {
	_, err := url.Parse(urlstr)
	if err != nil {
		return err
	}
	if runtime.GOOS == "linux" {
		_, err = exec.LookPath("xdg-open")
		if err != nil {
			return nil
		}
	}
	return browser.OpenURL(urlstr)
}
