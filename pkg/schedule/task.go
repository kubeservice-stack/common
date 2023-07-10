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
	"log"
	"reflect"
	"time"
)

// Task struct keeping information about job
type Task struct {
	interval uint64                   // pause interval * unit between runs
	taskFunc string                   // the job taskFunc to run, func[taskFunc]
	unit     timeUnit                 // time units, ,e.g. 'minutes', 'hours'...
	atTime   time.Duration            // optional time at which this job runs
	err      error                    // error related to job
	loc      *time.Location           // optional timezone that the atTime is in
	lastRun  time.Time                // datetime of last run
	nextRun  time.Time                // datetime of next run
	startDay time.Weekday             // Specific day of the week to start on
	funcs    map[string]interface{}   // Map for the function task store
	fparams  map[string][]interface{} // Map for function and  params of function
	lock     bool                     // lock the job from running at same time form multiple instances
	tags     []string                 // allow the user to tag jobs with certain labels
}

// NewTask creates a new job with the time interval.
func NewTask(interval uint64) *Task {
	return &Task{
		interval: interval,
		loc:      time.Local,
		lastRun:  time.Unix(0, 0),
		nextRun:  time.Unix(0, 0),
		startDay: time.Sunday,
		funcs:    make(map[string]interface{}),
		fparams:  make(map[string][]interface{}),
		tags:     []string{},
	}
}

// True if the job should be run now
func (j *Task) shouldRun() bool {
	return time.Now().Unix() >= j.nextRun.Unix()
}

// Run the job and immediately reschedule it
func (j *Task) run() ([]reflect.Value, error) {
	if j.lock {
		if locker == nil {
			return nil, fmt.Errorf("trying to lock %s with nil locker", j.taskFunc)
		}
		key := getFunctionKey(j.taskFunc)

		locker.Lock(key)
		defer locker.Unlock(key)
	}
	result, err := callTaskFuncWithParams(j.funcs[j.taskFunc], j.fparams[j.taskFunc])
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Err should be checked to ensure an error didn't occur creating the job
func (j *Task) Err() error {
	return j.err
}

// Do specifies the taskFunc that should be called every time the job runs
func (j *Task) Do(taskFun interface{}, params ...interface{}) error {
	if j.err != nil {
		return j.err
	}

	typ := reflect.TypeOf(taskFun)
	if typ.Kind() != reflect.Func {
		return ErrNotAFunction
	}
	fname := getFunctionName(taskFun)
	j.funcs[fname] = taskFun
	j.fparams[fname] = params
	j.taskFunc = fname

	now := time.Now().In(j.loc)
	if !j.nextRun.After(now) {
		j.scheduleNextRun()
	}

	return nil
}

// DoSafely does the same thing as Do, but logs unexpected panics, instead of unwinding them up the chain
// Deprecated: DoSafely exists due to historical compatibility and will be removed soon. Use Do instead
func (j *Task) DoSafely(taskFun interface{}, params ...interface{}) error {
	recoveryWrapperFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Internal panic occurred: %s", r)
			}
		}()

		_, _ = callTaskFuncWithParams(taskFun, params)
	}

	return j.Do(recoveryWrapperFunc)
}

// At schedules job at specific time of day
//
//	s.Every(1).Day().At("10:30:01").Do(task)
//	s.Every(1).Monday().At("10:30:01").Do(task)
func (j *Task) At(t string) *Task {
	hour, min, sec, err := formatTime(t)
	if err != nil {
		j.err = ErrTimeFormat
		return j
	}
	// save atTime start as duration from midnight
	j.atTime = time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
	return j
}

// GetAt returns the specific time of day the job will run at
//
//	s.Every(1).Day().At("10:30").GetAt() == "10:30"
func (j *Task) GetAt() string {
	return fmt.Sprintf("%1.2d:%2.2d", j.atTime/time.Hour, (j.atTime%time.Hour)/time.Minute)
}

// Loc sets the location for which to interpret "At"
//
//	s.Every(1).Day().At("10:30").Loc(time.UTC).Do(task)
func (j *Task) Loc(loc *time.Location) *Task {
	j.loc = loc
	return j
}

// Tag allows you to add labels to a job
// they don't impact the functionality of the job.
func (j *Task) Tag(t string, others ...string) {
	j.tags = append(j.tags, t)
	j.tags = append(j.tags, others...)
}

// Untag removes a tag from a job
func (j *Task) Untag(t string) {
	newTags := []string{}
	for _, tag := range j.tags {
		if t != tag {
			newTags = append(newTags, tag)
		}
	}

	j.tags = newTags
}

// Tags returns the tags attached to the job
func (j *Task) Tags() []string {
	return j.tags
}

func (j *Task) periodDuration() (time.Duration, error) {
	interval := time.Duration(j.interval)
	var periodDuration time.Duration

	switch j.unit {
	case seconds:
		periodDuration = interval * time.Second
	case minutes:
		periodDuration = interval * time.Minute
	case hours:
		periodDuration = interval * time.Hour
	case days:
		periodDuration = interval * time.Hour * 24
	case weeks:
		periodDuration = interval * time.Hour * 24 * 7
	default:
		return 0, ErrPeriodNotSpecified
	}
	return periodDuration, nil
}

// roundToMidnight truncate time to midnight
func (j *Task) roundToMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, j.loc)
}

// scheduleNextRun Compute the instant when this job should run next
func (j *Task) scheduleNextRun() error {
	now := time.Now()
	if j.lastRun == time.Unix(0, 0) {
		j.lastRun = now
	}

	periodDuration, err := j.periodDuration()
	if err != nil {
		return err
	}

	switch j.unit {
	case seconds, minutes, hours:
		j.nextRun = j.lastRun.Add(periodDuration)
	case days:
		j.nextRun = j.roundToMidnight(j.lastRun)
		j.nextRun = j.nextRun.Add(j.atTime)
	case weeks:
		j.nextRun = j.roundToMidnight(j.lastRun)
		dayDiff := int(j.startDay)
		dayDiff -= int(j.nextRun.Weekday())
		if dayDiff != 0 {
			j.nextRun = j.nextRun.Add(time.Duration(dayDiff) * 24 * time.Hour)
		}
		j.nextRun = j.nextRun.Add(j.atTime)
	}

	// advance to next possible schedule
	for j.nextRun.Before(now) || j.nextRun.Before(j.lastRun) {
		j.nextRun = j.nextRun.Add(periodDuration)
	}

	return nil
}

// NextScheduledTime returns the time of when this job is to run next
func (j *Task) NextScheduledTime() time.Time {
	return j.nextRun
}

// set the job's unit with seconds,minutes,hours...
func (j *Task) mustInterval(i uint64) error {
	if j.interval != i {
		return fmt.Errorf("interval must be %d", i)
	}
	return nil
}

// From schedules the next run of the job
func (j *Task) From(t *time.Time) *Task {
	j.nextRun = *t
	return j
}

// setUnit sets unit type
func (j *Task) setUnit(unit timeUnit) *Task {
	j.unit = unit
	return j
}

// Seconds set the unit with seconds
func (j *Task) Seconds() *Task {
	return j.setUnit(seconds)
}

// Minutes set the unit with minute
func (j *Task) Minutes() *Task {
	return j.setUnit(minutes)
}

// Hours set the unit with hours
func (j *Task) Hours() *Task {
	return j.setUnit(hours)
}

// Days set the job's unit with days
func (j *Task) Days() *Task {
	return j.setUnit(days)
}

// Weeks sets the units as weeks
func (j *Task) Weeks() *Task {
	return j.setUnit(weeks)
}

// Second sets the unit with second
func (j *Task) Second() *Task {
	j.mustInterval(1)
	return j.Seconds()
}

// Minute sets the unit  with minute, which interval is 1
func (j *Task) Minute() *Task {
	j.mustInterval(1)
	return j.Minutes()
}

// Hour sets the unit with hour, which interval is 1
func (j *Task) Hour() *Task {
	j.mustInterval(1)
	return j.Hours()
}

// Day sets the job's unit with day, which interval is 1
func (j *Task) Day() *Task {
	j.mustInterval(1)
	return j.Days()
}

// Week sets the job's unit with week, which interval is 1
func (j *Task) Week() *Task {
	j.mustInterval(1)
	return j.Weeks()
}

// Weekday start job on specific Weekday
func (j *Task) Weekday(startDay time.Weekday) *Task {
	j.mustInterval(1)
	j.startDay = startDay
	return j.Weeks()
}

// GetWeekday returns which day of the week the job will run on
// This should only be used when .Weekday(...) was called on the job.
func (j *Task) GetWeekday() time.Weekday {
	return j.startDay
}

// Monday set the start day with Monday
// - s.Every(1).Monday().Do(task)
func (j *Task) Monday() (job *Task) {
	return j.Weekday(time.Monday)
}

// Tuesday sets the job start day Tuesday
func (j *Task) Tuesday() *Task {
	return j.Weekday(time.Tuesday)
}

// Wednesday sets the job start day Wednesday
func (j *Task) Wednesday() *Task {
	return j.Weekday(time.Wednesday)
}

// Thursday sets the job start day Thursday
func (j *Task) Thursday() *Task {
	return j.Weekday(time.Thursday)
}

// Friday sets the job start day Friday
func (j *Task) Friday() *Task {
	return j.Weekday(time.Friday)
}

// Saturday sets the job start day Saturday
func (j *Task) Saturday() *Task {
	return j.Weekday(time.Saturday)
}

// Sunday sets the job start day Sunday
func (j *Task) Sunday() *Task {
	return j.Weekday(time.Sunday)
}

// Lock prevents job to run from multiple instances of gocron
func (j *Task) Lock() *Task {
	j.lock = true
	return j
}
