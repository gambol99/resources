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

package utils

import (
	"time"

	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	apiv1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	"github.com/gambol99/resources/pkg/client/clientset/versioned"
)

// UpdateCloudStatus is responsible for updating a cloud status
func UpdateCloudStatus(client versioned.Interface, status *apiv1.CloudStatus) error {
	return Retry(3, time.Second*2, func() error {
		// @check if the status already exists
		if _, err := client.CloudV1().CloudStatuses(status.Namespace).Get(status.Name, metav1.GetOptions{}); err != nil {
			if kerrors.IsNotFound(err) {
				_, err = client.CloudV1().CloudStatuses(status.Namespace).Create(status)
			}
			return err
		}
		_, err := client.CloudV1().CloudStatuses(status.Namespace).Update(status)

		return err
	})
}

// DeleteCloudStatus is responsible for updating a cloud status
func DeleteCloudStatus(client versioned.Interface, name, namespace string) error {
	return Retry(3, time.Second*2, func() error {
		// @check if the status already exists
		err := client.CloudV1().CloudStatuses(namespace).Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return nil
			}
			return err
		}

		return nil
	})
}

// UpdateKubernetesSecret is resposible for updating / creating a kube secret
func UpdateKubernetesSecret(client kubernetes.Interface, name, namespace string, data map[string]string) error {
	secret := &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		StringData: data,
	}

	return Retry(3, time.Duration(2*time.Second), func() error {
		if _, err := client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{}); err != nil {
			if kerrors.IsNotFound(err) {
				_, err = client.CoreV1().Secrets(namespace).Create(secret)
			}
			return err
		}
		_, err := client.CoreV1().Secrets(namespace).Update(secret)

		return err
	})
}

// FindKubernetesSecret is resposible for retrieving secrets from kubernetes
func FindKubernetesSecret(client kubernetes.Interface, name, namespace string) (map[string]string, error) {
	var err error
	var secret *core.Secret

	err = Retry(3, time.Duration(200*time.Millisecond), func() error {
		secret, err = client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
		return err
	})
	if err != nil {
		return map[string]string{}, err
	}

	return secret.StringData, nil
}

// FindCloudTemplate is responsible for retrieving the cloud template
func FindCloudTemplate(client versioned.Interface, name string) (*apiv1.CloudTemplate, error) {
	return client.Cloud().CloudTemplates().Get(name, metav1.GetOptions{})
}

// FindCloudResource is responsible for retrieving a cloud resource
func FindCloudResource(client versioned.Interface, name, namespace string) (*apiv1.CloudResource, error) {
	return client.Cloud().CloudResources(namespace).Get(name, metav1.GetOptions{})
}
