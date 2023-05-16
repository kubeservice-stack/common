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

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kubeservice-stack/common/pkg/schedule"
)

func ExampleTask(name string) {
	fmt.Println("Example Task " + name)

	t := time.NewTicker(time.Millisecond * 100)
	c := make(chan struct{})

	time.AfterFunc(time.Second*3, func() {
		close(c)
	})

	for {
		select {
		case <-t.C:
			fmt.Println(".")
		case <-c:
			fmt.Println("")
			return
		}
	}
}

type CustomLocker struct {
	Data map[string]interface{} //对非携程安全的数据，安全添加和删除
}

func (c *CustomLocker) Lock(key string) (bool, error) {
	c.Data[key] = time.Now()
	return true, nil
}

func (c *CustomLocker) Unlock(key string) error {
	delete(c.Data, key)
	return nil
}

func main() {
	var name string
	flag.StringVar(&name, "task-name", "example", "The example task name. Default: example")

	l := &CustomLocker{
		Data: make(map[string]interface{}),
	}

	schedule.SetLocker(l)

	schedule.Every(1).Second().Lock().Do(ExampleTask, name)
	<-schedule.Start()
}
