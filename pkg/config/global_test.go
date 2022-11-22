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

package config

import (
	"os"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
)

func TestEnvOverrides(t *testing.T) {
	var (
		envKey   string = "GIN_APP"
		expected string = "server-override"
	)
	os.Setenv(envKey, expected)
	defer func() { os.Unsetenv(envKey) }()

	GlobalCfg = Global{
		GinConfig: GinConfig{
			App:         "server",
			EnablePprof: true,
		},
	}
	err := env.Parse(&GlobalCfg)
	assert.Equal(t, err, nil)
	assert.Equal(t, GlobalCfg.GinConfig.App, expected)
	assert.Equal(t, GlobalCfg.GinConfig.EnablePprof, true)
}
