package ltoml

import (
	"fmt"

	"github.com/kubeservice-stack/common/pkg/utils"
)

func LoadConfig(cfgPath, defaultCfgPath string, v interface{}) error {
	if cfgPath == "" {
		cfgPath = defaultCfgPath
	}
	if !utils.Exist(cfgPath) {
		return fmt.Errorf("config file doesn't exist`")
	}

	if err := DecodeToml(cfgPath, v); err != nil {
		return fmt.Errorf("decode config file error:%s", err)
	}
	return nil
}
