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

package storage

import (
	"sort"

	"github.com/kubeservice-stack/common/pkg/utils"
)

const (
	// The maximum length of label name.
	// Longer names are truncated.
	maxLabelNameLen = 256

	// The maximum length of label value.
	// Longer values are truncated.
	maxLabelValueLen = 16 * 1024
)

// Label is a time-series label.
type Label struct {
	Key   string
	Value string
}

func marshalMetricName(metric string, labels []Label) string {
	if len(labels) == 0 {
		return metric
	}
	invalid := func(key, value string) bool {
		return key == "" || value == ""
	}

	// Determine the bytes size in advance.
	size := len(metric) + 2
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].Key < labels[j].Key
	})
	for i := range labels {
		label := &labels[i]
		if invalid(label.Key, label.Value) {
			continue
		}
		if len(label.Key) > maxLabelNameLen {
			label.Key = label.Key[:maxLabelNameLen]
		}
		if len(label.Value) > maxLabelValueLen {
			label.Value = label.Value[:maxLabelValueLen]
		}
		size += len(label.Key)
		size += len(label.Value)
		size += 4
	}

	// Start building the bytes.
	out := make([]byte, 0, size)
	out = utils.Uint16Encode(out, uint16(len(metric)))
	out = append(out, metric...)
	for i := range labels {
		label := &labels[i]
		if invalid(label.Key, label.Value) {
			continue
		}
		out = utils.Uint16Encode(out, uint16(len(label.Key)))
		out = append(out, label.Key...)
		out = utils.Uint16Encode(out, uint16(len(label.Value)))
		out = append(out, label.Value...)
	}
	return string(out)
}
