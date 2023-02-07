/*
Copyright 2023 The KubeService-Stack Authors.

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

package codec

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type Instance func() Codec

var adapters = make(map[PACK]Instance)

func Register(name PACK, adapter Instance) {
	if adapter == nil {
		panic("Codec: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("Codec: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func PluginInstance(name PACK) (adapter Codec) {
	instanceFunc, ok := adapters[name]
	if !ok {
		return
	}
	adapter = instanceFunc()
	return
}

func HasRegister(name PACK) bool {

	if _, ok := adapters[name]; ok {
		return true
	}
	return false
}
