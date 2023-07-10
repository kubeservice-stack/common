/*
Copyright 2022 The KubeService-Stack Authors.

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
	"github.com/mcuadros/go-version"
)

// Usage
//
//	Utils.CompareSimple("1.2", "1.0.1")
//	Returns: 1
//
//	Utils.CompareSimple("1.0rc1", "1.0")
//	Returns: -1
func Version_compare(version1, version2 string) int {
	return version.CompareSimple(version1, version2)
}
