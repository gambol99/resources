//go:generate ../../../../vendor/k8s.io/code-generator/generate-groups.sh all github.com/gambol99/resources/pkg/client github.com/gambol99/resources/pkg/apis resources:v1
// +k8s:deepcopy-gen=package,register

// Package v1 is the v1 version of the API.
// +groupName=cloud.appvia.io
package v1
