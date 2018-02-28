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

package aws

import (
	"context"
	"time"

	"github.com/gambol99/resources/pkg/models"
	log "github.com/sirupsen/logrus"
)

var (
	defaultWaitCheckInterval = time.Second * 5
)

// Wait is responsible for waiting for a stack to complete or fail
func (p *provider) Wait(ctx context.Context, name string, options *models.WaitOptions) (string, error) {
	interval := defaultWaitCheckInterval
	if options != nil && options.CheckInterval > 0 {
		interval = options.CheckInterval
	}

	ticker := time.NewTicker(interval)

	// @check the stack exists
	if found, err := p.hasStack(ctx, name); err != nil {
		return models.StatusUnknown, err
	} else if !found {
		return models.StatusUnknown, models.ErrStackNotFound
	}

	// @step: wait for check interval or signal to end
	started := time.Now()
	for {
		select {
		case <-ctx.Done():
			return models.StatusUnknown, models.ErrOperationAborted
		case <-ticker.C:

			status, err := p.Status(ctx, name, nil)
			if err != nil {
				return models.StatusUnknown, err
			}

			log.WithFields(log.Fields{
				"since":     time.Since(started).String(),
				"stackname": name,
				"status":    status,
			}).Debug("checking the status stack")

			if status == models.StatusInProgress {
				continue
			}

			return status, nil
		}
	}
}
