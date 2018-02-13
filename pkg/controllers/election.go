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
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"

	"github.com/gambol99/resources/pkg/controllers/api"
)

// election is a wrapper for an kubernetes election
type election struct {
	// client is the election service
	client *leaderelection.LeaderElector
}

// newElection creates a new election
func newElection(client kubernetes.Interface, recorder record.EventRecorder, name, namespace string) (api.Leadership, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("unable to get the hostname for election identity, error: %s", err)
	}

	lec, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock: &resourcelock.EndpointsLock{
			EndpointsMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      name,
			},
			Client: client.CoreV1(),
			LockConfig: resourcelock.ResourceLockConfig{
				Identity:      hostname,
				EventRecorder: recorder,
			},
		},
		LeaseDuration: 10 * time.Second,
		RenewDeadline: 8 * time.Second,
		RetryPeriod:   5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(<-chan struct{}) {},
			OnStoppedLeading: func() {},
			OnNewLeader:      func(leader string) {},
		},
	})
	if err != nil {
		return nil, err
	}

	// @step: start the election loop
	log.Info("starting the controller election service client")
	go lec.Run()

	return &election{client: lec}, nil
}

// IsLeader checks if the controller is a leader
func (e *election) IsLeader() bool {
	return e.client.IsLeader()
}
