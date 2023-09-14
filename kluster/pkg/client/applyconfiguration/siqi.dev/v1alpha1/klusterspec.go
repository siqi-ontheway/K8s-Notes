/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// KlusterSpecApplyConfiguration represents an declarative configuration of the KlusterSpec type for use
// with apply.
type KlusterSpecApplyConfiguration struct {
	Name        *string                      `json:"name,omitempty"`
	Region      *string                      `json:"region,omitempty"`
	Version     *string                      `json:"version,omitempty"`
	TokenSecret *string                      `json:"tokenSecret,omitempty"`
	NodePools   []NodePoolApplyConfiguration `json:"nodePools,omitempty"`
}

// KlusterSpecApplyConfiguration constructs an declarative configuration of the KlusterSpec type for use with
// apply.
func KlusterSpec() *KlusterSpecApplyConfiguration {
	return &KlusterSpecApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *KlusterSpecApplyConfiguration) WithName(value string) *KlusterSpecApplyConfiguration {
	b.Name = &value
	return b
}

// WithRegion sets the Region field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Region field is set to the value of the last call.
func (b *KlusterSpecApplyConfiguration) WithRegion(value string) *KlusterSpecApplyConfiguration {
	b.Region = &value
	return b
}

// WithVersion sets the Version field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Version field is set to the value of the last call.
func (b *KlusterSpecApplyConfiguration) WithVersion(value string) *KlusterSpecApplyConfiguration {
	b.Version = &value
	return b
}

// WithTokenSecret sets the TokenSecret field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TokenSecret field is set to the value of the last call.
func (b *KlusterSpecApplyConfiguration) WithTokenSecret(value string) *KlusterSpecApplyConfiguration {
	b.TokenSecret = &value
	return b
}

// WithNodePools adds the given value to the NodePools field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the NodePools field.
func (b *KlusterSpecApplyConfiguration) WithNodePools(values ...*NodePoolApplyConfiguration) *KlusterSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNodePools")
		}
		b.NodePools = append(b.NodePools, *values[i])
	}
	return b
}