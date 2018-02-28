/*
Copyright 2018 All rights reserved - Appvia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version

import (
	"fmt"
	"strconv"
	"time"
)

var (
	// Release is the application version
	Release = "v0.0.1"
	// GitSHA is the commit we were built off
	GitSHA = "unknown"
	// BuildTime is the build time
	BuildTime = "0"
	// Author is the application author
	Author = "Rohith Jayawardene"
	// Email is the author email
	Email = "rohith.jayawardene@appvia.io"
)

// GetVersion returns the version tag
func GetVersion() string {
	return fmt.Sprintf("%s (git+sha %s)", Release, GitSHA)
}

// GetBuildTime returns the build time of the application
func GetBuildTime() time.Time {
	tm, err := strconv.ParseInt(BuildTime, 10, 64)
	if err != nil {
		panic("unable to parse the build time")
	}

	return time.Unix(tm, 0)
}
