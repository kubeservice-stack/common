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
