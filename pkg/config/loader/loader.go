// configloader

/*
配置加载后, 联动更新各个模块的全局变量.
*/
package configloader

import (
	"github.com/kubeservice-stack/common/pkg/config"
	"github.com/kubeservice-stack/common/pkg/logger"
	"github.com/kubeservice-stack/common/pkg/metrics"
)

func LoadConfig(cfgPath string) (err error) {
	if err = config.LoadGlobalConfig(cfgPath); err != nil {
		return
	}
	// 更新各个模块的全局变量
	// Logging
	if err = logger.NewLogger(config.GlobalCfg.Logging); err != nil {
		return
	}
	return
}
