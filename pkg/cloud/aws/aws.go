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
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	log "github.com/sirupsen/logrus"

	"github.com/gambol99/resources/pkg/models"
	"github.com/gambol99/resources/pkg/utils"
)

// provider implements the aws provider
type provider struct {
	// the iam client
	accounts iamiface.IAMAPI
	// the cloud formation client
	client cloudformationiface.CloudFormationAPI
	// the compute client
	compute ec2iface.EC2API
	// the configuration for the provider
	config *models.ProviderConfig
}

const (
	// maxAccessKeys is the max number of keys a user can have
	maxAccessKeys = 2
)

// New creates and returns a aws provider
func New(config *models.ProviderConfig) (models.CloudProvider, error) {
	if config.Region == "" {
		log.Debug("no aws region has been specified, using the metadata service or environment variables")
		config.Region = findRegion()
	}
	// @check we have a region
	if config.Region == "" {
		return nil, errors.New("you must specifiy the aws region, no metedata service available")
	}
	if config.ClusterName == "" {
		return nil, errors.New("you have not set the clustername")
	}
	if config.Name == "" {
		return nil, errors.New("you have not set the provider name")
	}

	sess, err := session.NewSession(&aws.Config{
		MaxRetries: aws.Int(5),
		Region:     aws.String(config.Region),
	})
	if err != nil {
		return nil, err
	}

	return &provider{
		accounts: iam.New(sess),
		client:   cloudformation.New(sess),
		compute:  ec2.New(sess),
		config:   config,
	}, nil
}

// findRegion attempts to find the region from the metadata service
func findRegion() string {
	var region string

	// @step: check the environment variable
	for _, x := range []string{"AWS_REGION", "AWS_DEFAULT_REGION"} {
		if region = os.Getenv(x); region != "" {
			log.Debugf("using the aws region: %s", region)
			return region
		}
	}

	// @step: attempt to pull from the EC2 metadata service
	client := ec2metadata.New(session.Must(session.NewSession()))
	if !client.Available() {
		return region
	}

	utils.Retry(5, time.Second*5, func() error {
		resp, err := client.GetInstanceIdentityDocument()
		if err != nil {
			return fmt.Errorf("failed to retrieve instance document from metadata service: %s", err)
		}
		region = resp.Region

		return nil
	})

	return region
}
