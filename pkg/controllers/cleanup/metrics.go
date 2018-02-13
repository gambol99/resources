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

package cleanup

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	cleanupCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cleanup_run_total",
			Help: "The total number of invocations for the cleanup controller",
		},
	)
	cleanupDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "cleanup_duration_seconds",
			Help: "Cleanup latency distributions.",
		},
	)
	cleanupErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cleanup_error_total",
			Help: "A total number of errors encountered by the cleanup controller",
		},
	)
)

func init() {
	prometheus.MustRegister(cleanupCounter)
	prometheus.MustRegister(cleanupDuration)
	prometheus.MustRegister(cleanupErrors)
}
