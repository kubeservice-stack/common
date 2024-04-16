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
	"testing"

	"github.com/kubeservice-stack/common/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_MySQL(t *testing.T) {
	assert := assert.New(t)
	r := NewMySQL(config.GlobalCfg.DBConfig)
	assert.NotNil(r)
}

func TestString2Map(t *testing.T) {
	assert := assert.New(t)
	ret := String2Map("")
	assert.Equal(ret, map[string]string{})

	ret = String2Map("  asdfsdf ?")
	assert.Equal(ret, map[string]string{})

	ret = String2Map("  dd=dd ?")
	assert.Equal(ret, map[string]string{"dd": "dd"})

	ret = String2Map("  dd=dd=dd ?")
	assert.Equal(ret, map[string]string{})

	ret = String2Map("  dd=dd&dd=aa ?")
	assert.Equal(ret, map[string]string{"dd": "aa"})

	ret = String2Map("  dd=dd&dd=aa&bb=aa ?")
	assert.Equal(ret, map[string]string{"dd": "aa", "bb": "aa"})
}
