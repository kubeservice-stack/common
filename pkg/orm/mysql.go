/*
Copyright 2024 The KubeService-Stack Authors.

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

package orm

import (
	"strings"
	"time"

	driver_mysql "github.com/go-sql-driver/mysql"
	"github.com/kubeservice-stack/common/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(cfg config.DBConfig) gorm.Dialector {
	loc, _ := time.LoadLocation("UTC")
	dcfg := driver_mysql.Config{
		User:                 cfg.User,
		Passwd:               cfg.Password,
		Net:                  "tcp",
		Addr:                 cfg.Host + ":" + cfg.Port,
		DBName:               cfg.Database,
		Params:               String2Map(cfg.Args),
		Loc:                  loc,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	return mysql.Open(dcfg.FormatDSN())
}

// [?param1=value1&...&paramN=valueN]
func String2Map(str string) map[string]string {
	str1 := strings.Trim(str, "?")
	v := strings.Split(str1, "&")
	ret := make(map[string]string)
	for _, value := range v {
		kv := strings.Split(strings.Trim(value, " "), "=")
		if len(kv) == 2 {
			ret[kv[0]] = kv[1]
		}
	}
	return ret
}

func init() {
	Register(config.MYSQL, NewMySQL)
}
