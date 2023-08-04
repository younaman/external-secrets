/*
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

package v1alpha1

import (
	"context"

	"github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Ready indicates that the client is configured correctly
	// and can be used.
	ValidationResultReady ValidationResult = iota

	// Unknown indicates that the client can be used
	// but information is missing and it can not be validated.
	ValidationResultUnknown

	// Error indicates that there is a misconfiguration.
	ValidationResultError
)

type ValidationResult uint8

func (v ValidationResult) String() string {
	return [...]string{"Ready", "Unknown", "Error"}[v]
}

// +kubebuilder:object:root=false
// +kubebuilder:object:generate:false
// +k8s:deepcopy-gen:interfaces=nil
// +k8s:deepcopy-gen=nil

// Provider is a common interface for interacting with secret backends.
type Provider interface {
	// NewClient constructs a SecretsManager Provider
	NewClientV2(ctx context.Context, config runtime.Object, kube client.Client, namespace string) (v1beta1.SecretsClient, error)
}
