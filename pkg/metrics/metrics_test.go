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

package metrics

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kubeservice-stack/common/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
)

func Test_DefaultRegistry(t *testing.T) {
	assert := assert.New(t)
	r := DefaultRegistry()
	assert.NotNil(r)
}

func Test_CustomTallyScopeConfig(t *testing.T) {
	assert := assert.New(t)
	DefaultTallyScope = NewTallyScope(CustomTallyScopeConfig("aaaa"))
	counter := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "bar",
	}).Counter("test_counter")

	gauge := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "baz",
	}).Gauge("test_gauge")

	timer := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "qux",
	}).Timer("test_timer_summary")

	histogram := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "quk",
	}).Histogram("test_histogram", tally.DefaultBuckets)

	go func() {
		for {
			counter.Inc(1)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			gauge.Update(rand.Float64() * 1000)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			tsw := timer.Start()
			hsw := histogram.Start()
			time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
			tsw.Stop()
			hsw.Stop()
		}
	}()

	time.Sleep(1 * time.Second)
	router := gin.New()

	router.Any("/metrics", gin.WrapH(DefaultTallyScope.Reporter.HTTPHandler()))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(http.StatusOK, w.Code)

	log.Println(w.Body.String())

	err := DefaultTallyScope.Destroy()
	assert.Nil(err)
}

func Test_DefaultTallyScope_close(t *testing.T) {
	assert := assert.New(t)
	aa := &config.Metrics{
		FlushInterval:          5,
		EnableGoRuntimeMetrics: false,
		MetricsPrefix:          "aa_aa",
		MetricsTags:            map[string]string{"aa": "test"},
	}
	DefaultTallyScope = NewTallyScope(aa)

	counter := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "bar",
	}).Counter("test_counter")

	gauge := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "baz",
	}).Gauge("test_gauge")

	timer := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "qux",
	}).Timer("test_timer_summary")

	histogram := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "quk",
	}).Histogram("test_histogram", tally.DefaultBuckets)

	go func() {
		for {
			counter.Inc(1)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			gauge.Update(rand.Float64() * 1000)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			tsw := timer.Start()
			hsw := histogram.Start()
			time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
			tsw.Stop()
			hsw.Stop()
		}
	}()

	time.Sleep(1 * time.Second)
	router := gin.New()

	router.Any("/metrics", gin.WrapH(DefaultTallyScope.Reporter.HTTPHandler()))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(http.StatusOK, w.Code)

	log.Println(w.Body.String())

	err := DefaultTallyScope.Destroy()
	assert.Nil(err)
}

func Test_DefaultTallyScope(t *testing.T) {
	assert := assert.New(t)
	DefaultTallyScope = NewTallyScope(&config.GlobalCfg.Metrics)

	counter := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "bar",
	}).Counter("test_counter")

	counteraa := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "baaar",
	}).Counter("test_counter")

	aabb := DefaultTallyScope.Scope.SubScope("aa").Tagged(map[string]string{
		"foo": "baaar",
	}).Counter("test_counter")

	gauge := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "baz",
	}).Gauge("test_gauge")

	timer := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foaao": "aaa",
	}).Timer("test_timer_summary")

	histogram := DefaultTallyScope.Scope.Tagged(map[string]string{
		"foo": "quk",
	}).Histogram("test_histogram", tally.DefaultBuckets)

	go func() {
		for {
			counter.Inc(1)
			counteraa.Inc(1)
			aabb.Inc(2)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			gauge.Update(rand.Float64() * 1000)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			tsw := timer.Start()
			hsw := histogram.Start()
			time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
			tsw.Stop()
			hsw.Stop()
		}
	}()

	time.Sleep(1 * time.Second)
	router := gin.New()

	router.Any("/metrics", gin.WrapH(DefaultTallyScope.Reporter.HTTPHandler()))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(http.StatusOK, w.Code)

	log.Println(w.Body.String())

	err := DefaultTallyScope.Destroy()
	assert.Nil(err)
}
