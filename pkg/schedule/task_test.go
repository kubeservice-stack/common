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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func CallBackTest(ss interface{}) interface{} {
	return ss
}

func CallBackPanic(ss interface{}) {
	panic("error panic")
}

func Test_NewTask(t *testing.T) {
	assert := assert.New(t)
	task := NewTask(1)
	now := time.Now()
	err := task.At(now.Format("15:04:05")).Do(CallBackTest, "aaa")
	assert.Nil(err)
	err = task.At(now.Format("15:04:05")).Do(CallBackTest)
	assert.Nil(err)

	ddd := task.GetAt()
	assert.Equal(ddd, now.Format("15:04"))

	eee := task.GetWeekday()
	assert.Equal(eee, time.Weekday(0))

	err = task.Days().Day().Friday().Weeks().Weekday(time.Monday).Wednesday().Tuesday().Thursday().
		Sunday().Saturday().Monday().At(time.Now().Add(1*time.Second).Format("15:04:05")).Do(CallBackTest, "aaa")
	assert.Nil(err)

	var functionNot interface{} = nil

	err = task.Hour().Minutes().Minute().Hours().Second().Seconds().Do("ddd", "df")
	assert.Equal(err, ErrNotAFunction)

	err = task.Hour().Minutes().Minute().Hours().Second().Seconds().DoSafely(functionNot, "df")
	assert.Nil(err)
}

func Test_TaskTags(t *testing.T) {
	assert := assert.New(t)
	task := NewTask(1)
	now := time.Now()
	err := task.Week().At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackTest, "aaa")
	assert.Nil(err)

	task.Tag("dd", "dff", "ddd", "dd")
	aa := task.Tags()
	assert.Equal(aa, []string{"dd", "dff", "ddd", "dd"})

	task.Untag("aa")
	aa = task.Tags()
	assert.Equal(aa, []string{"dd", "dff", "ddd", "dd"})
	task.Untag("dd")
	aa = task.Tags()
	assert.Equal(aa, []string{"dff", "ddd"})

	err = task.From(&now).At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackTest, "aaa")
	assert.Nil(err)

	task1 := NewTask(1)
	err = task1.Days().At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackTest, "aaa")
	assert.Nil(err)

	aeea := task1.NextScheduledTime()
	assert.NotNil(aeea.Format("2006-01-02 15:04:05"))

	err = task1.Minute().At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackTest, "aaa")
	assert.Nil(err)

	err = task1.Second().At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackTest, "aaa")
	assert.Nil(err)

	err = task1.Hour().At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackTest, "aaa")
	assert.Nil(err)

	err = task1.Day().At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackTest, "aaa")
	assert.Nil(err)

	err = task1.Day().At(now.Format("24:04:05")).Loc(time.UTC).Do(CallBackTest, "aaa")
	assert.NotNil(err)
}

func Test_CallPanic(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()

	task1 := NewTask(1)
	err := task1.Seconds().At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackPanic)
	time.Sleep(2 * time.Second)
	assert.Nil(err)

	err = task1.Seconds().At(now.Format("15:04:05")).Loc(time.UTC).DoSafely(CallBackPanic)
	time.Sleep(2 * time.Second)
	assert.Equal(err, task1.Err())

	err = task1.Seconds().Lock().At(now.Format("15:04:05")).Loc(time.UTC).DoSafely(CallBackPanic)
	assert.Equal(err, task1.Err())
}

type fakeLocker struct{}

func (s *fakeLocker) Lock(key string) (bool, error) {
	return true, nil
}

func (s *fakeLocker) Unlock(key string) error {
	return nil
}

func Test_TaskLock(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	task1 := NewTask(1)
	SetLocker(&fakeLocker{})
	err := task1.Lock().Minutes().At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackPanic)
	assert.Nil(err)

	err = task1.Lock().Hours().At(now.Format("15:04:05")).Loc(time.UTC).Do(CallBackPanic)
	assert.Nil(err)
}
