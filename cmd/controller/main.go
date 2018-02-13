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

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"

	"github.com/gambol99/resources/pkg/controllers"
	"github.com/gambol99/resources/pkg/controllers/api"
	"github.com/gambol99/resources/pkg/version"
)

func main() {
	app := cli.NewApp()
	app.Usage = "provides a kubernetes controller for the automated management of cloud resources via custom resources"
	app.Version = version.GetVersion()
	app.Author = version.Author
	app.Email = version.Email
	app.Compiled = version.GetBuildTime()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "name",
			Usage:  "the name for this controller, used when creating the stacks `NAME`",
			EnvVar: "PROVIDER_NAME",
			Value:  "resource.appvia.io/default",
		},
		cli.BoolTFlag{
			Name:   "enable-metrics",
			Usage:  "indicated you wish to enable the metrics endpoint `BOOL`",
			EnvVar: "ENABLE_METRICS",
		},
		cli.StringFlag{
			Name:   "cloud-provider",
			Usage:  "the cloud provider implemetation, i.e. aws, gce etc `NAME`",
			EnvVar: "CLOUD_PROVIDER",
			Value:  "aws",
		},
		cli.DurationFlag{
			Name:   "resync-duration",
			Usage:  "the time duration for a force resync of controller state `DURATION`",
			EnvVar: "RESYNC_DURATION",
			Value:  0,
		},
		cli.DurationFlag{
			Name:   "stack-timeout",
			Usage:  "the default timeout for a stack tom complete or error `DURATION`",
			EnvVar: "STACK_TIMEOUT",
			Value:  time.Minute * 30,
		},
		cli.StringFlag{
			Name:   "election-namespace",
			Usage:  "the namespace for used for the controller election `NAMESPACE`",
			EnvVar: "KUBE_NAMESPACE",
			Value:  "kube-system",
		},
		cli.IntFlag{
			Name:   "threadness",
			Usage:  "the number of worker routines for the controller worker `NUMBER`",
			EnvVar: "THREADNESS",
			Value:  1,
		},
		cli.StringFlag{
			Name:   "metrics-listen",
			Usage:  "the interface we should expose the controller metrics on `INTERFACE`",
			EnvVar: "METRICS_LISTEN",
			Value:  "127.0.0.1:8080",
		},
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "indicates verbose logging on the controller `BOOL`",
			EnvVar: "VERBOSE",
		},
	}
	app.Action = func(cx *cli.Context) error {
		// @step: create the controller
		return func() error {
			c, err := controllers.New(&api.Config{
				CloudProvider:     cx.String("cloud-provider"),
				EnableMetrics:     cx.Bool("enable-metrics"),
				ElectionNamespace: cx.String("election-namespace"),
				MetricsListen:     cx.String("metrics-listen"),
				Name:              cx.String("name"),
				ResyncDuration:    cx.Duration("resync-duration"),
				StackTimeout:      cx.Duration("stack-timeout"),
				Threadness:        cx.Int("threadness"),
				Verbose:           cx.Bool("verbose"),
			})
			if err != nil {
				return err
			}

			// @step: provide a upper context for the controller to run under
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				if err := c.Run(ctx); err != nil {
					fmt.Fprintf(os.Stderr, "failed to start controller: %s", err)
					os.Exit(1)
				}
			}()

			signalChannel := make(chan os.Signal)
			signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			<-signalChannel

			// @step: wait for the resource controller to gracefully shutdown
			return c.Wait(time.Duration(10 * time.Minute))
		}()
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err)
		os.Exit(1)
	}
}
