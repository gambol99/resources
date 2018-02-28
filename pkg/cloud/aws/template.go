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
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gambol99/resources/pkg/models"
)

// Templater providers a cloudformation templater
type Templater struct {
	ctx    context.Context
	client ec2iface.EC2API
	config *models.ProviderConfig
}

// Render is responsibe for generating the template
func (t *Templater) Render(c context.Context, values map[string]string, content string) (string, error) {
	var err error

	t.ctx = c
	tm := template.New("main")
	if _, err = tm.Funcs(t.templateFuncsMap(tm)).Parse(content); err != nil {
		return "", err
	}
	tm.Option("missingkey=error")

	// @step: render the actual template
	writer := new(bytes.Buffer)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to parse template, error: %s", r)
		}
	}()

	// @step; used to capture the time of a render
	capture := prometheus.NewTimer(templateDuration)
	defer capture.ObserveDuration()

	if err = tm.ExecuteTemplate(writer, "name", values); err != nil {
		return "", err
	}

	return writer.String(), nil
}

// NewTemplater creates and returns a templater
func NewTemplater(client ec2iface.EC2API, config *models.ProviderConfig) *Templater {
	return &Templater{
		client: client,
		config: config,
	}
}

// Region returns the region we are in
func (t *Templater) Region() string {
	return t.config.Region
}

// Filter is responsible for filtering on objects
func (t *Templater) Filter(items []models.Object) []models.Object {
	var selected []models.Object

	return selected
}

// Subnets returns a list of subnets within the VPC
func (t *Templater) Subnets() []models.Network {
	resp, err := t.client.DescribeSubnetsWithContext(t.ctx, &ec2.DescribeSubnetsInput{
		Filters: getClusterFilters(t.config.ClusterName),
	})
	if err != nil {
		t.panic("failed to describe the subnets: %s", err)
	}

	var list []models.Network
	for _, x := range resp.Subnets {
		tags := makeTags(x.Tags)
		name := tags["Name"]
		list = append(list, models.Network{
			AvailabilityZone: aws.StringValue(x.AvailabilityZone),
			CIDR:             aws.StringValue(x.CidrBlock),
			Object: models.Object{
				ID:   aws.StringValue(x.SubnetId),
				Name: name,
				Tags: tags,
			},
		})
	}

	return list
}

// Network returns a list of networks
func (t *Templater) Network() models.Network {
	resp, err := t.client.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: getClusterFilters(t.config.ClusterName),
	})
	if err != nil {
		t.panic("failed to describe the vpcs, error: %s", err)
	}
	if len(resp.Vpcs) <= 0 {
		t.panic("failed to find any vpcs")
	}
	if len(resp.Vpcs) > 1 {
		t.panic("found more than one vpc")
	}
	tags := makeTags(resp.Vpcs[0].Tags)

	return models.Network{
		CIDR: aws.StringValue(resp.Vpcs[0].CidrBlock),
		Object: models.Object{
			ID:   aws.StringValue(resp.Vpcs[0].VpcId),
			Name: tags["Name"],
			Tags: tags,
		},
	}
}

// NetworkID return the vpc id
func (t *Templater) NetworkID() string {
	return t.Network().ID
}

// templateFuncsMap returns a map if the template functions for this template
func (t *Templater) templateFuncsMap(tm *template.Template) template.FuncMap {
	//funcs := sprig.TxtFuncMap()
	funcs := make(map[string]interface{}, 0)
	funcs["region"] = t.Region
	funcs["vpc"] = t.Network
	funcs["vpcid"] = t.NetworkID
	funcs["filter"] = t.Filter
	funcs["subnets"] = t.Subnets

	return funcs
}

// panic is responsible for throwing a panic
func (t *Templater) panic(message string, opts ...interface{}) {
	panic(fmt.Sprintf(message, opts...))
}

func makeTags(tags []*ec2.Tag) map[string]string {
	list := make(map[string]string, 0)
	for _, x := range tags {
		list[aws.StringValue(x.Key)] = aws.StringValue(x.Value)
	}

	return list
}
