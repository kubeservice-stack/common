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
	"sort"
	"time"
)

// Scheduler struct, the only data member is the list of jobs.
type Scheduler struct {
	jobs [MAXJOBNUM]*Task // Array store jobs
	size int              // Size of jobs which jobs holding.
	loc  *time.Location   // Location to use when scheduling jobs with specified times
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: [MAXJOBNUM]*Task{},
		size: 0,
		loc:  time.Local,
	}
}

// Tasks returns the list of Tasks from the Scheduler
func (s *Scheduler) Tasks() []*Task {
	return s.jobs[:s.size]
}

func (s *Scheduler) Len() int {
	return s.size
}

func (s *Scheduler) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

func (s *Scheduler) Less(i, j int) bool {
	return s.jobs[j].nextRun.Unix() >= s.jobs[i].nextRun.Unix()
}

// ChangeLoc changes the default time location
func (s *Scheduler) ChangeLoc(newLocation *time.Location) {
	s.loc = newLocation
}

// Get the current runnable jobs, which shouldRun is True
func (s *Scheduler) getRunnableTasks() (runningTasks [MAXJOBNUM]*Task, n int) {
	runnableTasks := [MAXJOBNUM]*Task{}
	n = 0
	sort.Sort(s)
	for i := 0; i < s.size; i++ {
		if s.jobs[i].shouldRun() {
			runnableTasks[n] = s.jobs[i]
			n++
		} else {
			break
		}
	}
	return runnableTasks, n
}

// NextRun datetime when the next job should run.
func (s *Scheduler) NextRun() (*Task, time.Time) {
	if s.size <= 0 {
		return nil, time.Now()
	}
	sort.Sort(s)
	return s.jobs[0], s.jobs[0].nextRun
}

// Every schedule a new periodic job with interval
func (s *Scheduler) Every(interval uint64) *Task {
	job := NewTask(interval).Loc(s.loc)
	s.jobs[s.size] = job
	s.size++
	return job
}

// RunPending runs all the jobs that are scheduled to run.
func (s *Scheduler) RunPending() {
	runnableTasks, n := s.getRunnableTasks()

	if n != 0 {
		for i := 0; i < n; i++ {
			go runnableTasks[i].run()
			runnableTasks[i].lastRun = time.Now()
			runnableTasks[i].scheduleNextRun()
		}
	}
}

// RunAll run all jobs regardless if they are scheduled to run or not
func (s *Scheduler) RunAll() {
	s.RunAllwithDelay(0)
}

// RunAllwithDelay runs all jobs with delay seconds
func (s *Scheduler) RunAllwithDelay(d int) {
	for i := 0; i < s.size; i++ {
		go s.jobs[i].run()
		if d != 0 {
			time.Sleep(time.Duration(d))
		}
	}
}

// Remove specific job j by function
func (s *Scheduler) Remove(j interface{}) {
	s.removeByCondition(func(someTask *Task) bool {
		return someTask.taskFunc == getFunctionName(j)
	})
}

// RemoveByRef removes specific job j by reference
func (s *Scheduler) RemoveByRef(j *Task) {
	s.removeByCondition(func(someTask *Task) bool {
		return someTask == j
	})
}

// RemoveByTag removes specific job j by tag
func (s *Scheduler) RemoveByTag(t string) {
	s.removeByCondition(func(someTask *Task) bool {
		for _, a := range someTask.tags {
			if a == t {
				return true
			}
		}
		return false
	})
}

func (s *Scheduler) removeByCondition(shouldRemove func(*Task) bool) {
	i := 0

	// keep deleting until no more jobs match the criteria
	for {
		found := false

		for ; i < s.size; i++ {
			if shouldRemove(s.jobs[i]) {
				found = true
				break
			}
		}

		if !found {
			return
		}

		for j := (i + 1); j < s.size; j++ {
			s.jobs[i] = s.jobs[j]
			i++
		}
		s.size--
		s.jobs[s.size] = nil
	}
}

// Scheduled checks if specific job j was already added
func (s *Scheduler) Scheduled(j interface{}) bool {
	for _, job := range s.jobs {
		if job.taskFunc == getFunctionName(j) {
			return true
		}
	}
	return false
}

// Clear delete all scheduled jobs
func (s *Scheduler) Clear() {
	for i := 0; i < s.size; i++ {
		s.jobs[i] = nil
	}
	s.size = 0
}

// Start all the pending jobs
// Add seconds ticker
func (s *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.RunPending()
			case <-stopped:
				ticker.Stop()
				return
			}
		}
	}()

	return stopped
}

// The following methods are shortcuts for not having to
// create a Scheduler instance

// Every schedules a new periodic job running in specific interval
func Every(interval uint64) *Task {
	return defaultScheduler.Every(interval)
}

// RunPending run all jobs that are scheduled to run
//
// Please note that it is *intended behavior that run_pending()
// does not run missed jobs*. For example, if you've registered a job
// that should run every minute and you only call run_pending()
// in one hour increments then your job won't be run 60 times in
// between but only once.
func RunPending() {
	defaultScheduler.RunPending()
}

// RunAll run all jobs regardless if they are scheduled to run or not.
func RunAll() {
	defaultScheduler.RunAll()
}

// RunAllwithDelay run all the jobs with a delay in seconds
//
// A delay of `delay` seconds is added between each job. This can help
// to distribute the system load generated by the jobs more evenly over
// time.
func RunAllwithDelay(d int) {
	defaultScheduler.RunAllwithDelay(d)
}

// Start run all jobs that are scheduled to run
func Start() chan bool {
	return defaultScheduler.Start()
}

// Clear all scheduled jobs
func Clear() {
	defaultScheduler.Clear()
}

// Remove a specific job
func Remove(j interface{}) {
	defaultScheduler.Remove(j)
}

// Scheduled checks if specific job j was already added
func Scheduled(j interface{}) bool {
	for _, job := range defaultScheduler.jobs {
		if job.taskFunc == getFunctionName(j) {
			return true
		}
	}
	return false
}

// NextRun gets the next running time
func NextRun() (job *Task, time time.Time) {
	return defaultScheduler.NextRun()
}
