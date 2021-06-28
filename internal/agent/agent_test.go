package agent

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/michael-diggin/proglog/internal/config"
	"github.com/michael-diggin/proglog/internal/loadbalance"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	api "github.com/michael-diggin/proglog/api/v1"
)

func TestAgent(t *testing.T) {
	serverTLSConfig, err := config.SetUpTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	peerTLSConfig, err := config.SetUpTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		Server:        false,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	var agents []*Agent
	for i := 0; i < 3; i++ {
		bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", getFreePort())
		rpcPort := getFreePort()
		dataDir, err := ioutil.TempDir("", "agent-test-log")
		require.NoError(t, err)

		var startJoinAddrs []string
		if i != 0 {
			startJoinAddrs = []string{agents[0].Config.BindAddr}
		}
		agent, err := New(Config{
			NodeName:        fmt.Sprintf("%d", i),
			StartJoinAddrs:  startJoinAddrs,
			BindAddr:        bindAddr,
			RPCPort:         rpcPort,
			DataDir:         dataDir,
			ACLModelFile:    config.ACLModelFile,
			ACLPolicyFile:   config.ACLPolicyFile,
			ServerTLSConfig: serverTLSConfig,
			PeerTLSConfig:   peerTLSConfig,
			Bootstrap:       i == 0,
		})
		require.NoError(t, err)
		agents = append(agents, agent)
	}
	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(agent.Config.DataDir))
		}
	}()
	time.Sleep(1 * time.Second)

	message := &api.Record{Value: []byte("hello world")}
	ctx := context.Background()
	leaderClient := client(t, agents[0], peerTLSConfig)
	produceResponse, err := leaderClient.Produce(ctx, &api.ProduceRequest{Record: message})
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	consumeResponse, err := leaderClient.Consume(ctx, &api.ConsumeRequest{Offset: produceResponse.Offset})
	require.NoError(t, err)
	require.Equal(t, message.Value, consumeResponse.Record.Value)

	followClient := client(t, agents[1], peerTLSConfig)
	consumeResponse, err = followClient.Consume(ctx, &api.ConsumeRequest{Offset: produceResponse.Offset})
	require.NoError(t, err)
	require.Equal(t, message.Value, consumeResponse.Record.Value)

	consumeResponse, err = leaderClient.Consume(ctx, &api.ConsumeRequest{Offset: produceResponse.Offset + 1})
	require.Error(t, err)
	require.Nil(t, consumeResponse)
	require.Equal(t, codes.NotFound, grpc.Code(err))
}

func client(t *testing.T, agent *Agent, tlsConfig *tls.Config) api.LogClient {
	tlsCreds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}
	rpcAddr, err := agent.Config.RPCAddr()
	require.NoError(t, err)
	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", loadbalance.Name, rpcAddr), opts...)
	require.NoError(t, err)
	return api.NewLogClient(conn)
}

func getFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}
