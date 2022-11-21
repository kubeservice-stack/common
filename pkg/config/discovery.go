package config

import (
	"encoding/json"
	"fmt"

	"github.com/kubeservice-stack/common/pkg/utils"
)

type Discovery struct {
	Namespace   string         `toml:"namespace" json:"namespace" env:"DISCOVERY_NAMESPACE"`         // 命名空间
	Endpoints   []string       `toml:"endpoints" json:"endpoints" env:"DISCOVERY_ENDPOINTS"`         // 连接端点
	DialTimeout utils.Duration `toml:"dial_timeout" json:"dial_timeout" env:"DISCOVERY_DIALTIMEOUT"` // 连接超时时间
}

func (ds Discovery) TOML() string {
	if len(ds.Endpoints) == 0 {
		ds.Endpoints = []string{}
	}
	endpoints, _ := json.Marshal(ds.Endpoints)
	return fmt.Sprintf(`
[discovery]
  ## etcd namespace
  namespace = "%s"
  ## etcd 集群配置
  endpoints = %s
  ## ETCD连接 timeout时间
  dial_timeout = "%s"`,
		ds.Namespace,
		endpoints,
		ds.DialTimeout.String(),
	)
}

func (ds Discovery) DefaultConfig() Discovery {
	ds = Discovery{
		Namespace: "application",
		Endpoints: []string{"http://127.0.0.1:2379"},
	}
	return ds
}
