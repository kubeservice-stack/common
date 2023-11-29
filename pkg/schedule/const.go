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
	"errors"
	"time"
)

const MAXJOBNUM = 1000

var (
	ErrTimeFormat           = errors.New("time format error")
	ErrParamsNotAdapted     = errors.New("the number of params is not adapted")
	ErrNotAFunction         = errors.New("only functions can be schedule into the job queue")
	ErrPeriodNotSpecified   = errors.New("unspecified job period")
	ErrParameterCannotBeNil = errors.New("nil parameters cannot be used with reflection")
)

type timeUnit int

const (
	seconds timeUnit = iota + 1
	minutes
	hours
	days
	weeks
)

func (t timeUnit) String() time.Duration {
	switch t {
	case seconds:
		return time.Second
	case minutes:
		return time.Minute
	case hours:
		return time.Hour
	case days:
		return time.Hour * 24
	case weeks:
		return time.Hour * 24 * 7
	default:
		return -1
	}
}

var (
	locker           Locker
	defaultScheduler = NewScheduler()
)
