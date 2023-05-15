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

package schedule

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Schdule(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	sched := NewScheduler()
	err := sched.Every(1).At(now.Format("15:04:05")).DoSafely(fmt.Println, "aa")
	assert.Nil(err)
	sched.RunAll()
	go func() {
		<-sched.Start()
	}()
	time.Sleep(1 * time.Second)
	sched.Clear()
}

func Test_MutiSchdule(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	sched := NewScheduler()
	err := sched.Every(1).At(now.Format("15:04:05")).DoSafely(fmt.Println, "aa")
	assert.Nil(err)
	err = sched.Every(1).At(now.Format("15:04:05")).DoSafely(fmt.Println, "bb")
	assert.Nil(err)
	sched.RunAll()

	l := sched.Len()
	assert.Equal(l, 2)

	tasks := sched.Tasks()
	assert.Equal(len(tasks), 2)

	sched.ChangeLoc(time.UTC)
	assert.Equal(sched.loc.String(), "UTC")

	task, tm := sched.NextRun()
	assert.NotEmpty(tm)
	assert.NotNil(task)

	sched.Remove(fmt.Println)
	sched.Remove(fmt.Printf)
	sched.Remove(fmt.Errorf)
	sched.RemoveByTag("aa")

	err = sched.Every(1).At(now.Format("15:04:05")).DoSafely(fmt.Println, "cc")
	assert.Nil(err)

	go func() {
		<-sched.Start()
	}()
	time.Sleep(1 * time.Second)
	sched.Clear()
}

func Test_DefaultScheduler(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()

	err := Every(1).At(now.Format("15:04:05")).DoSafely(fmt.Println, "aa")
	assert.Nil(err)
	err = Every(1).At(now.Format("15:04:05")).DoSafely(fmt.Println, "bb")
	assert.Nil(err)
	RunAll()
	RunAllwithDelay(1)
	RunPending()

	task, tm := NextRun()
	assert.NotEmpty(tm)
	assert.NotNil(task)

	go func() {
		<-Start()
	}()
	time.Sleep(1 * time.Second)
	Remove(fmt.Append)
	Clear()
}
