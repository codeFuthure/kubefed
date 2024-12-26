/*
Copyright 2024 The CodeFuture Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FederatedResource is a generic representation of a federated type
type FederatedResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	ClusterStatus []ResourceClusterStatus `json:"clusterStatus,omitempty"`
}

// ResourceClusterStatus defines the status of federated resource within a cluster
type ResourceClusterStatus struct {
	ClusterName string                 `json:"clusterName,omitempty"`
	Status      map[string]interface{} `json:"status,omitempty"`
}