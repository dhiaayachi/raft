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

type NotImplemented struct {
}

func (e NotImplemented) Error() string {
	return fmt.Sprintf("Not Implemented")
}

type GRPCTransport struct {
	Server   *grpc.Server
	cp       *grpcpool.Pool
	addr     raft.ServerAddress
	consumer chan raft.RPC
}

func (g GRPCTransport) SetHeartbeatHandler(cb func(rpc raft.RPC)) {
}

func (g GRPCTransport) Consumer() <-chan raft.RPC {
	return g.consumer
}

func (g GRPCTransport) LocalAddr() raft.ServerAddress {
	return g.addr
}

func (g GRPCTransport) AppendEntriesPipeline(id raft.ServerID, target raft.ServerAddress) (raft.AppendPipeline, error) {
	return nil, NotImplemented{}
}

func (g GRPCTransport) AppendEntries(id raft.ServerID, target raft.ServerAddress, args *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) error {
	ctx := context.Background()
	fmt.Printf("dhayachi:: grpc appendentry\n")
	r := new(commands.AppendEntriesRequest)
	commands.AppendEntriesRequestFromStruct(args, r)
	conn, err := grpc.Dial(string(target), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		fmt.Printf("dhayachi:: falling back to append entry 1\n")
		return err
	}
	client := commands.NewCommandsClient(conn)
	entries, err := client.AppendEntries(ctx, r)
	if err != nil {
		fmt.Printf("dhayachi:: falling back to append entry 2\n")
		return err
	}
	commands.AppendEntriesResponseToStruct(entries, resp)

	return nil
}

func (g GRPCTransport) RequestVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestVoteRequest, resp *raft.RequestVoteResponse) error {
	return NotImplemented{}
}

func (g GRPCTransport) InstallSnapshot(id raft.ServerID, target raft.ServerAddress, args *raft.InstallSnapshotRequest, resp *raft.InstallSnapshotResponse, data io.Reader) error {
	return NotImplemented{}
}

func (g GRPCTransport) EncodePeer(id raft.ServerID, addr raft.ServerAddress) []byte {
	return []byte(addr)
}

func (g GRPCTransport) DecodePeer(bytes []byte) raft.ServerAddress {
	return raft.ServerAddress(bytes)
}

func (g GRPCTransport) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	return NotImplemented{}
}

// NewGRPCTransport is used to initialize a new transport and
// generates a random local address if none is specified. The given timeout
// will be used to decide how long to wait for a connected peer to process the
// RPCs that we're sending it. See also Connect() and Consumer().
func NewGRPCTransport(addr raft.ServerAddress, timeout time.Duration) (raft.ServerAddress, raft.Transport) {
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
	rpcs := make(chan raft.RPC)
	srv := &commands.AppendEntriesServerService{RPCch: rpcs}
	commands.RegisterCommandsServer(server, srv)
	return addr, &GRPCTransport{Server: server, addr: addr, consumer: rpcs}

}
