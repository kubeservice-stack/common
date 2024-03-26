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

package config

import (
	"fmt"
)

type DBTYPE string

const (
	MYSQL    DBTYPE = "mysql"
	SQLITE3  DBTYPE = "sqlite3"
	POSTGRES DBTYPE = "postgres"
)

// DBConfg represents a database configuration
type DBConfg struct {
	DBType   DBTYPE `toml:"db_type" json:"db_type" env:"ORM_DBTYPE"`     // dbtype
	User     string `toml:"user" json:"user" env:"ORM_USER"`             // user
	Password string `toml:"password" json:"password" env:"ORM_PASSWORD"` // password
	Host     string `toml:"host" json:"host" env:"ORM_HOST"`             // host
	Port     string `toml:"port" json:"port" env:"ORM_PORT"`             // port
	Database string `toml:"database" json:"database" env:"ORM_DATABASE"` // database
	CaCert   string `toml:"ca_cert" json:"ca_cert" env:"ORM_CACERT"`     // cacert
}

func (db DBConfg) TOML() string {
	return fmt.Sprintf(`
[database]
  ## dbtype 模式: mysql/sqlite3/postgres
  db_type = "%v"
  ## user 用户
  user = "%s"
  ## password 密码
  password = "%s"
  ## Host 地址
  host = "%s"
  ## Port 端口
  port = "%s"
  ## 数据库名
  database = "%s"
  ## ca cert证书
  ca_cert = "%s"`,
		string(db.DBType),
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Database,
		db.CaCert)
}

func (db DBConfg) DefaultConfig() DBConfg {
	db = DBConfg{
		DBType:   MYSQL,
		User:     "root",
		Password: "root",
		Host:     "localhost",
		Port:     "3306",
		Database: "test",
	}
	return db
}
