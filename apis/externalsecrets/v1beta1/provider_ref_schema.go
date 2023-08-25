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

package v1beta1

import (
	"fmt"
	"sync"

	esmeta "github.com/external-secrets/external-secrets/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var refBuilder map[string]Provider
var refbuildlock sync.RWMutex

func init() {
	refBuilder = make(map[string]Provider)
}

// Register a store backend type. Register panics if a
// backend with the same store is already registered.
func RefRegister(s Provider, storeType string) {
	buildlock.Lock()
	defer buildlock.Unlock()
	_, exists := builder[storeType]
	if exists {
		panic(fmt.Sprintf("store %q already registered", storeType))
	}

	builder[storeType] = s
}

// ForceRegister adds to store schema, overwriting a store if
// already registered. Should only be used for testing.
func RefForceRegister(s Provider, storeType string) {
	buildlock.Lock()
	builder[storeType] = s
	buildlock.Unlock()
}

// GetProviderByName returns the provider implementation by name.
func GetProviderByRef(ref esmeta.ProviderRef) (Provider, bool) {
	buildlock.RLock()
	f, ok := builder[ref.Kind]
	buildlock.RUnlock()
	return f, ok
}

// GetProvider returns the provider from the generic store.
func GetProviderByObj(r runtime.Object) (Provider, error) {
	if r == nil {
		return nil, nil
	}
	kind := r.GetObjectKind().GroupVersionKind().Kind
	if kind == "" {
		return nil, fmt.Errorf("no spec found in %#v", r)
	}
	buildlock.RLock()
	f, ok := builder[kind]
	buildlock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("failed to find registered store backend for type: %s", kind)
	}

	return f, nil
}
