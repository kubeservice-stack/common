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
