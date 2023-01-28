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

package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/kubeservice-stack/common/pkg/config"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
)

func Test_Logger(t *testing.T) {
	logger1 := GetLogger("pkg/common/logger", "test")
	RunningAtomicLevel.SetLevel(zapcore.DebugLevel)

	logger1.Warn("warn for test", String("count", "1"), Reflect("v1", map[string]string{"a": "1"}))
	logger1.Info("info for test", Uint16("value", 1), Int32("v1", 2),
		Int64("v2", 2), Any("v3", map[string]string{"a": "1"}))
	logger1.Debug("debug for test", Uint32("value", 2))
	logger1.Error("error for test", Error(fmt.Errorf("error")))

	assert.NotNil(t, defaultLogger)

	logger3 := GetLogger("pkg/common/logger", "")
	logger3.Error("error test")
}

func Test_Access_logger(t *testing.T) {
	assert.Nil(t, NewLogger(config.Logging{Level: "debug"}))
	logger1 := GetLogger(HTTPModule, "access")
	logger1.Info("access log")
	isTerminal = true
	defer func() {
		isTerminal = false
	}()
	logger1.Info("access log")
}

func Test_Level_String(t *testing.T) {
	isTerminal = true
	defer func() {
		isTerminal = false
	}()
	assert.Equal(t, "DEBUG", LevelString(zapcore.DebugLevel))
	assert.Equal(t, "DPANIC", LevelString(zapcore.DPanicLevel))
	assert.Equal(t, "INFO", LevelString(zapcore.InfoLevel))
	assert.Equal(t, "WARN", LevelString(zapcore.WarnLevel))
	assert.Equal(t, "ERROR", LevelString(zapcore.ErrorLevel))
	isTerminal = false
	assert.Equal(t, "ERROR", LevelString(zapcore.ErrorLevel))
}

func Test_Logger_Stack(t *testing.T) {
	panicFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				GetLogger(CrashModule, "test-painc").
					getInitializedOrDefaultLogger().Panic("panic stack", Stack())
			}
		}()
		panic("test-panic")
	}
	assert.Panics(t, panicFunc)
}

func Test_IsTerminal(t *testing.T) {
	assert.False(t, IsTerminal(os.Stdout))
}

func Test_InitLogger(t *testing.T) {
	assert.NotNil(t, GetLogger("test", "test").getInitializedOrDefaultLogger())

	cfg1 := config.Logging{Level: "LLL"}
	assert.NotNil(t, NewLogger(cfg1))

	assert.Nil(t, NewLogger(config.GlobalCfg.Logging))
	thisLogger := GetLogger("test", "test")
	assert.NotNil(t, thisLogger.getInitializedOrDefaultLogger())
	assert.NotNil(t, thisLogger.getInitializedOrDefaultLogger())

	cfg3 := config.Logging{Level: "info"}
	assert.Nil(t, NewLogger(cfg3))

	cfg4 := config.Logging{Level: "debug"}
	assert.Nil(t, NewLogger(cfg4))

	isTerminal = true
	defer func() {
		isTerminal = false
	}()
	assert.Nil(t, NewLogger(cfg4))
}
