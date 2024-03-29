//go:build !aix && !windows
// +build !aix,!windows

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

package lock

import (
	"fmt"
	"os"
	"syscall"
)

// Lock try locking file, return err if fails.
func (l *fileLock) lock() error {
	// invoke syscall for file lock
	if err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return fmt.Errorf("cannot flock directory %s - %s", l.fileName, err)
	}
	return nil
}

// Unlock unlock file lock, if fail return err
func (l *fileLock) unlock() error {
	return syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
}

func (l *fileLock) trylock() bool {
	fileInfo, err := os.Stat(l.fileName)
	if err == nil && fileInfo != nil {
		return false
	}
	return true
}
