package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"

	"github.com/kubeservice-stack/common/pkg/config/ltoml"
)

var GlobalCfg Global

func init() {
	GlobalCfg.Logging = GlobalCfg.Logging.DefaultConfig()
	GlobalCfg.Metrics = GlobalCfg.Metrics.DefaultConfig()
	GlobalCfg.Discovery = GlobalCfg.Discovery.DefaultConfig()
	GlobalCfg.GinConfig = GlobalCfg.GinConfig.DefaultConfig()
	GlobalCfg.RateLimit = GlobalCfg.RateLimit.DefaultConfig()
	GlobalCfg.Temporary = GlobalCfg.Temporary.DefaultConfig()
}

func LoadGlobalConfig(cfgPath string) (err error) {
	// 文件配置覆盖默认值
	if err = ltoml.LoadConfig(cfgPath, cfgPath, &GlobalCfg); err != nil {
		return
	}
	// 环境变量覆盖配置的值
	if err = env.Parse(&GlobalCfg); err != nil {
		return
	}
	return
}

type Global struct {
	Logging   `toml:"logging"`
	Metrics   `toml:"metrics"`
	Discovery `toml:"discovery"`
	GinConfig `toml:"gin"`
	RateLimit `toml:"ratelimit"`
	Temporary `toml:"temporary"`
}

func (g Global) TOML() string {
	return fmt.Sprint(
		g.Logging.TOML(),
		g.Metrics.TOML(),
		g.Discovery.TOML(),
		g.GinConfig.TOML(),
		g.RateLimit.TOML(),
		g.Temporary.TOML(),
	)
}
