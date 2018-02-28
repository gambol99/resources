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

package null

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gambol99/resources/pkg/models"
)

type provider struct {
	sync.RWMutex

	stacks map[string]*models.Stack
}

// New returns a null provider
func New(config *models.ProviderConfig) (models.CloudProvider, error) {
	log.Info("creating a new null cloud provider")
	return &provider{
		stacks: make(map[string]*models.Stack, 0),
	}, nil
}

// Credentials generates the credentials from a stack
func (p *provider) Credentials(context.Context, string) ([]models.Credential, error) {
	return []models.Credential{}, nil
}

// Create is responsible for creating or updating a stack
func (p *provider) Create(ctx context.Context, name string, options *models.CreateOptions) error {
	log.WithFields(log.Fields{
		"name":      name,
		"resource":  options.Resource.Name,
		"namespace": options.Resource.Namespace,
		"template":  options.Template.Name,
		"model":     options.Context,
	}).Info("creating a new stack")

	resource := options.Resource

	stack := &models.Stack{
		Created:   time.Now(),
		Name:      name,
		Namespace: resource.Namespace,
		Spec: models.StackSpec{
			Name:      resource.Name,
			Retention: time.Duration(time.Hour * 24),
			Tags:      options.Tags,
			Template:  options.Template.Name,
		},
		Status: models.StackStatus{
			Status: models.StatusDone,
		},
	}

	p.Lock()
	defer p.Unlock()

	p.stacks[name] = stack

	return nil
}

// Logs returns the logs from a stack
func (p *provider) Logs(ctx context.Context, name string, options *models.GetOptions) (string, error) {
	_, err := p.getStack(ctx, name)

	return "", err
}

// Delete is responsible for removing the stack
func (p *provider) Delete(ctx context.Context, name string, options *models.DeleteOptions) error {
	log.WithFields(log.Fields{
		"name": name,
	}).Info("deleting the stack")

	p.RLock()
	defer p.RUnlock()

	if _, found := p.stacks[name]; !found {
		return models.ErrStackNotFound
	}
	delete(p.stacks, name)

	return nil
}

// Exists is responsible for checking is stack already exists
func (p *provider) Exists(ctx context.Context, name string) (*models.Stack, bool, error) {
	stack, err := p.getStack(ctx, name)
	if err != nil {
		return nil, false, err
	}

	return stack, true, nil
}

// Get is responisble for retrieving a stack
func (p *provider) Get(ctx context.Context, name string, options *models.GetOptions) (*models.Stack, error) {
	return p.getStack(ctx, name)
}

// List is responsible for getting a list of stacks
func (p *provider) List(context.Context, *models.ListOptions) ([]*models.Stack, error) {
	var list []*models.Stack
	p.RLock()
	defer p.RUnlock()

	for _, x := range p.stacks {
		list = append(list, x)
	}

	return list, nil
}

// Status is responsible for getting the status
func (p *provider) Status(ctx context.Context, name string, options *models.GetOptions) (string, error) {
	stack, err := p.getStack(ctx, name)
	if err != nil {
		return models.StatusUnknown, err
	}

	return stack.Status.Status, nil
}

// UpdateTags is responsible for updating just the tags of a stack
func (p *provider) UpdateTags(ctx context.Context, name string, tags map[string]string) error {
	stack, err := p.getStack(ctx, name)
	if err != nil {
		return err
	}
	p.Lock()
	defer p.Unlock()

	for k, v := range tags {
		stack.Spec.Tags[k] = v
	}

	return nil
}

// Wait is responsible for waiting for a stack to complete or fail
func (p *provider) Wait(context.Context, string, *models.WaitOptions) (string, error) {
	return models.StatusDone, nil
}

func (p *provider) getStack(ctx context.Context, name string) (*models.Stack, error) {
	p.RLock()
	defer p.RUnlock()

	stack, found := p.stacks[name]
	if !found {
		return nil, models.ErrStackNotFound
	}

	return stack, nil
}
