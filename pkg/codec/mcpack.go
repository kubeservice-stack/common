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

import (
	"github.com/kubeservice-stack/common/pkg/codec/mcpack"
)

type MCPack struct{}

func NewMCPack() Codec {
	return &MCPack{}
}

func (mc *MCPack) Marshal(v interface{}) ([]byte, error) {
	return mcpack.Marshal(v)
}

func (mc *MCPack) Unmarshal(data []byte, v interface{}) error {
	return mcpack.Unmarshal(data, v)
}

func init() {
	Register(MCPACK, NewMCPack)
}
