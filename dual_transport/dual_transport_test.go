package dual_transport

import (
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft/grpc_transport"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGrpc(t *testing.T) {
	conf1 := raft.InmemConfig(t)
	conf1.LocalID = "grpc1"
	store1 := raft.NewInmemStore()
	snap1 := raft.NewInmemSnapshotStore()
	_, trans1 := raft.NewInmemTransport("127.0.0.1:12345")
	_, transport1 := grpc_transport.NewGRPCTransport("127.0.0.1:12345", time.Second)
	_, dt1 := NewDualTransport(transport1, trans1)
	r1, err := raft.NewRaft(conf1, &raft.MockFSM{}, store1, store1, snap1, dt1)
	require.NoError(t, err)
	require.NotNil(t, r1)

	conf2 := raft.InmemConfig(t)
	conf2.LocalID = "grpc2"
	store2 := raft.NewInmemStore()
	snap2 := raft.NewInmemSnapshotStore()
	_, trans2 := raft.NewInmemTransport("127.0.0.1:12346")
	_, transport2 := grpc_transport.NewGRPCTransport("127.0.0.1:12346", time.Second)
	_, dt2 := NewDualTransport(transport2, trans2)
	_, err = raft.NewRaft(conf2, &raft.MockFSM{}, store2, store2, snap2, dt2)

	require.NoError(t, err)

	conf3 := raft.InmemConfig(t)
	conf3.LocalID = "grpc3"
	store3 := raft.NewInmemStore()
	snap3 := raft.NewInmemSnapshotStore()
	_, trans3 := raft.NewInmemTransport("127.0.0.1:12347")
	_, transport3 := grpc_transport.NewGRPCTransport("127.0.0.1:12347", time.Second)
	_, dt3 := NewDualTransport(transport3, trans3)
	_, err = raft.NewRaft(conf3, &raft.MockFSM{}, store3, store3, snap3, dt3)

	require.NoError(t, err)

	trans1.Connect("127.0.0.1:12346", trans2)
	trans1.Connect("127.0.0.1:12347", trans3)
	r1.BootstrapCluster(raft.Configuration{Servers: []raft.Server{{ID: "grpc1", Address: "127.0.0.1:12345"}, {ID: "grpc2", Address: "127.0.0.1:12346"}, {ID: "grpc3", Address: "127.0.0.1:12347"}}})
	time.Sleep(10 * time.Second)
}
