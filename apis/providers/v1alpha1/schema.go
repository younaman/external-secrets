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
	"fmt"
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var refBuilder map[string]client.Object
var refBuildlock sync.RWMutex

func init() {
	refBuilder = make(map[string]client.Object)
}

func RefRegister(s client.Object, providerKind string) {
	refBuildlock.Lock()
	defer refBuildlock.Unlock()
	_, exists := refBuilder[providerKind]
	if exists {
		panic(fmt.Sprintf("provider %q already registered", providerKind))
	}

	refBuilder[providerKind] = s
}

// ForceRegister adds to store schema, overwriting a store if
// already registered. Should only be used for testing.
func RefForceRegister(s client.Object, storeType string) {
	refBuildlock.Lock()
	refBuilder[storeType] = s
	refBuildlock.Unlock()
}

// GetProvider returns the provider from the generic store.
func GetManifestByKind(kind string) (client.Object, error) {
	refBuildlock.RLock()
	f, ok := refBuilder[kind]
	refBuildlock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("failed to find registered object backend for kind: %s", kind)
	}
	return f, nil
}
