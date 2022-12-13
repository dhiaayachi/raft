package grpc_transport

import (
	"context"
	"fmt"
	"github.com/hashicorp/raft"
	commands "github.com/hashicorp/raft/proto"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"time"
)

type GrpcTransport struct {
	Server   *grpc.Server
	cp       *grpcpool.Pool
	addr     raft.ServerAddress
	consumer chan raft.RPC
	trans    raft.Transport
}

func (g GrpcTransport) Consumer() <-chan raft.RPC {
	return g.consumer
}

func (g GrpcTransport) LocalAddr() raft.ServerAddress {
	return g.addr
}

func (g GrpcTransport) AppendEntriesPipeline(id raft.ServerID, target raft.ServerAddress) (raft.AppendPipeline, error) {
	return g.trans.AppendEntriesPipeline(id, target)
}

func (g GrpcTransport) AppendEntries(id raft.ServerID, target raft.ServerAddress, args *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) error {
	ctx := context.Background()
	fmt.Printf("dhayachi:: grpc appendentry\n")
	r := new(commands.AppendEntriesRequest)
	commands.AppendEntriesRequestFromStruct(args, r)
	conn, err := grpc.Dial(string(target), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		fmt.Printf("dhayachi:: falling back to append entry 1\n")
		return g.trans.AppendEntries(id, target, args, resp)
	}
	client := commands.NewCommandsClient(conn)
	entries, err := client.AppendEntries(ctx, r)
	if err != nil {
		fmt.Printf("dhayachi:: falling back to append entry 2\n")
		return g.trans.AppendEntries(id, target, args, resp)
	}
	commands.AppendEntriesResponseToStruct(entries, resp)

	return nil
}

func (g GrpcTransport) RequestVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestVoteRequest, resp *raft.RequestVoteResponse) error {
	return g.trans.RequestVote(id, target, args, resp)
}

func (g GrpcTransport) InstallSnapshot(id raft.ServerID, target raft.ServerAddress, args *raft.InstallSnapshotRequest, resp *raft.InstallSnapshotResponse, data io.Reader) error {
	return g.trans.InstallSnapshot(id, target, args, resp, data)
}

func (g GrpcTransport) EncodePeer(id raft.ServerID, addr raft.ServerAddress) []byte {
	return g.trans.EncodePeer(id, addr)
}

func (g GrpcTransport) DecodePeer(bytes []byte) raft.ServerAddress {
	return g.trans.DecodePeer(bytes)
}

func (g GrpcTransport) SetHeartbeatHandler(cb func(rpc raft.RPC)) {
	g.trans.SetHeartbeatHandler(cb)
}

func (g GrpcTransport) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	return g.trans.TimeoutNow(id, target, args, resp)
}

// NewGrpcTransport is used to initialize a new transport and
// generates a random local address if none is specified. The given timeout
// will be used to decide how long to wait for a connected peer to process the
// RPCs that we're sending it. See also Connect() and Consumer().
func NewGrpcTransport(addr raft.ServerAddress, timeout time.Duration, trans raft.Transport) (raft.ServerAddress, raft.Transport) {
	server := grpc.NewServer(grpc.ConnectionTimeout(timeout))

	go func() {
		listen, err := net.Listen("tcp", string(addr))
		if err != nil {

			panic("listener creation error" + err.Error())
		}
		err = server.Serve(listen)
		if err != nil {
			panic("listener Serve error" + err.Error())
		}
	}()
	c := trans.Consumer()
	rpcs := make(chan raft.RPC)
	go func() {
		for {
			rpc := <-c
			rpcs <- rpc
		}
	}()
	srv := &commands.AppendEntriesServerService{RPCch: rpcs}
	commands.RegisterCommandsServer(server, srv)
	return addr, &GrpcTransport{Server: server, addr: addr, consumer: rpcs, trans: trans}

}
