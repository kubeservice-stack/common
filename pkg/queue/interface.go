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

package queue

import (
	"fmt"
)

var (
	ErrExceedingMessageSizeLimit = fmt.Errorf("message exceeds the max page size limit")
	ErrOutOfSequenceRange        = fmt.Errorf("out of sequence range")
	ErrExceedingTotalSizeLimit   = fmt.Errorf("queue data size exceeds the max size limit")
	ErrMsgNotFound               = fmt.Errorf("message not found")
)

type Queue interface {
	Push(item interface{})                     // queue 结尾put 数据， 如果失败返回err
	Pop() (interface{}, bool)                  // 获取消息数据
	Length() int64                             // queue 长度
	IsEmpty() bool                             // queue是否为空
	PopMany(count int64) ([]interface{}, bool) // 获取多条消息数据
}
