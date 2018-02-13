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

package resources

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricErrorTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "resource_controller_errors_total",
			Help: "The total number of errors encountered by the resource controller",
		},
	)
)

func init() {
	prometheus.MustRegister(metricErrorTotal)
}
