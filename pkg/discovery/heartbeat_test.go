package discovery

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	etcdcliv3 "go.etcd.io/etcd/clientv3"
)

type EtcdHeartbeatTestSuite struct {
	suite.Suite
	Cluster *ETCDMockCluster
}

func (s *EtcdHeartbeatTestSuite) SetupTest() {
	s.Cluster = StartEtcdMockCluster(s.T())
}

func (s *EtcdHeartbeatTestSuite) TearDownTest() {
	s.Cluster.Terminate(s.T())
}

func (s *EtcdHeartbeatTestSuite) TestHeartBeatKeepaliveStop() {
	cfg := etcdcliv3.Config{
		Endpoints: s.Cluster.Endpoints,
	}
	cli, err := etcdcliv3.New(cfg)
	s.Nil(err)

	heartbeat := newHeartbeat(cli, "/test/heartbeat", []byte("dongjiang"), 0, false)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ok, err := heartbeat.grantKeepAliveLease(ctx)
	s.Nil(err)
	s.Equal(true, ok)

	go func() {
		heartbeat.keepAlive(ctx)
	}()

	val, err := cli.Get(ctx, "/test/heartbeat")
	s.Nil(err)
	s.Equal(1, len(val.Kvs))

	err = cli.Close()
	s.Nil(err)
	time.Sleep(time.Second)
}

// go test 入口
func TestEtcdHeartbeatTestSuite(t *testing.T) {
	suite.Run(t, new(EtcdHeartbeatTestSuite))
}
