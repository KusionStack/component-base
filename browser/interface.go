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

// Browser is an object that knows how to open a given URL in a new tab in
// some suitable browser on the current system.
type Browser interface {
	// OpenURL opens the given URL in a web browser.
	//
	// Depending on the circumstances and on the target platform, this may or
	// may not cause the browser to take input focus. Because of this
	// uncertainty, any caller of this method must be sure to include some
	// language in its UI output to let the user know that a browser tab has
	// opened somewhere, so that they can go and find it if the focus didn't
	// switch automatically.
	OpenURL(url string) error
}
