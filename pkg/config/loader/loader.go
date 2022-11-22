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

/*
配置加载后, 联动更新各个模块的全局变量.
*/
package configloader

import (
	"github.com/kubeservice-stack/common/pkg/config"
	"github.com/kubeservice-stack/common/pkg/logger"
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
