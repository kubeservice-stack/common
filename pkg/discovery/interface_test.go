package discovery

import (
	"testing"

	"github.com/kubeservice-stack/common/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestNewRepo(t *testing.T) {
	assert := assert.New(t)
	cluster := StartEtcdMockCluster(t)
	defer cluster.Terminate(t)
	cfg := config.Discovery{
		Endpoints: cluster.Endpoints,
	}

	factory := NewDiscoveryFactory("nobody")
	ds, err := factory.CreateDiscovery(cfg)
	assert.Nil(err)
	assert.NotNil(ds)
}

func TestEventType_String(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("delete", EventTypeDelete.String())
	assert.Equal("modify", EventTypeModify.String())
	assert.Equal("all", EventTypeAll.String())
	assert.Equal("unknown", EventType(111).String())
}
