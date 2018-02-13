/*
Copyright 2018 All rights reserved - Appvia.io

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

package controllers

import (
	"context"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/gambol99/resources/pkg/client/clientset/versioned/fake"
)

func newTestController(t *testing.T) *resourceController {
	c := &resourceController{clientset: fake.NewSimpleClientset()}
	return c
}

func newTestRunningController(t *testing.T) (*resourceController, context.CancelFunc) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	c := newTestController(t)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := c.Run(ctx); err != nil {
			t.Fatalf("failed to start test controller, error: %s", err)
		}
	}()

	return c, cancel
}

func TestNewController(t *testing.T) {
	c := newTestController(t)
	assert.NotNil(t, c)
}

func TestRunController(t *testing.T) {
	c, cancel := newTestRunningController(t)
	defer cancel()
	assert.NotNil(t, c)
	assert.NotNil(t, cancel)
}
