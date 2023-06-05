package grpc_transport

import (
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGrpc(t *testing.T) {
	conf1 := raft.InmemConfig(t)
	conf1.LocalID = "grpc1"
	store1 := raft.NewInmemStore()
	snap1 := raft.NewInmemSnapshotStore()
	_, transport1 := NewGRPCTransport("127.0.0.1:12345", time.Second)
	r1, err := raft.NewRaft(conf1, &raft.MockFSM{}, store1, store1, snap1, transport1)
	require.NoError(t, err)
	require.NotNil(t, r1)

	conf2 := raft.InmemConfig(t)
	conf2.LocalID = "grpc2"
	store2 := raft.NewInmemStore()
	snap2 := raft.NewInmemSnapshotStore()
	_, transport2 := NewGRPCTransport("127.0.0.1:12346", time.Second)
	_, err = raft.NewRaft(conf2, &raft.MockFSM{}, store2, store2, snap2, transport2)

	require.NoError(t, err)

	conf3 := raft.InmemConfig(t)
	conf3.LocalID = "grpc3"
	store3 := raft.NewInmemStore()
	snap3 := raft.NewInmemSnapshotStore()
	_, transport3 := NewGRPCTransport("127.0.0.1:12347", time.Second)
	_, err = raft.NewRaft(conf3, &raft.MockFSM{}, store3, store3, snap3, transport3)

	require.NoError(t, err)

	r1.BootstrapCluster(raft.Configuration{Servers: []raft.Server{{ID: "grpc1", Address: "127.0.0.1:12345"}, {ID: "grpc2", Address: "127.0.0.1:12346"}, {ID: "grpc3", Address: "127.0.0.1:12347"}}})
	time.Sleep(10 * time.Second)
}
