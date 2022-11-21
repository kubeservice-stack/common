package config

import (
	"fmt"
)

type RateLimit struct {
	QPS   int `toml:"qps" json:"qps" env:"GIN_RATELIMIT_QPS"`       // qps
	Burst int `toml:"burst" json:"burst" env:"GIN_RATELIMIT_BURST"` // 并发数
}

func (rl RateLimit) TOML() string {
	return fmt.Sprintf(`
# 访问频率限制
[ratelimit]
  ## qps
  qps = %d
  ## 并发数
  burst = %d`,
		rl.QPS, rl.Burst,
	)
}

func (rl RateLimit) DefaultConfig() RateLimit {
	return RateLimit{
		QPS:   100,
		Burst: 20,
	}
}
