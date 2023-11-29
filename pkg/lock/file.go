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

	"github.com/kubeservice-stack/common/pkg/logger"
)

var openFileFn = os.OpenFile

// fileLock is file lock
type fileLock struct {
	fileName string
	file     *os.File

	logger *logger.Logger
}

// NewFileLock create new file lock instance
func NewFileLock(fileName string) (Locker, error) {
	f, err := openFileFn(fileName, os.O_CREATE|os.O_RDONLY, os.FileMode(0o600))
	if err != nil {
		return nil, fmt.Errorf("cannot create file[%s] for lock err: %s", fileName, err)
	}
	return &fileLock{
		file:     f,
		fileName: fileName,
		logger:   logger.GetLogger("Lockers", "FileLock"),
	}, nil
}

// Lock try locking file, return err if fails.
func (l *fileLock) Lock() error {
	return l.lock()
}

// Unlock unlock file lock, if fail return err
func (l *fileLock) Unlock() error {
	defer func() {
		if err := os.Remove(l.fileName); err != nil {
			l.logger.Error("remove file lock error", logger.String("file", l.fileName), logger.Error(err))
		}
		l.logger.Info("remove file lock successfully", logger.String("file", l.fileName))
	}()

	defer func() {
		if err := l.file.Close(); err != nil {
			l.logger.Error("close file lock error", logger.String("file", l.fileName), logger.Error(err))
		}
	}()
	return l.unlock()
}

//  TryLock will try to lock file and return whether it succeed or not without blocking.

func (l *fileLock) TryLock() bool {
	return l.trylock()
}
