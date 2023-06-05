package appendentries

import (
	"context"
	"github.com/hashicorp/raft"
)

// AppendEntriesServerService must be embedded to have forward compatible implementations.
type AppendEntriesServerService struct {
	UnimplementedCommandsServer
	RPCch chan raft.RPC
}

func (cs *AppendEntriesServerService) AppendEntries(ctx context.Context, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	r := new(raft.AppendEntriesRequest)
	AppendEntriesRequestToStruct(req, r)
	rpc := raft.RPC{}
	rpc.Command = r
	responses := make(chan raft.RPCResponse)
	rpc.RespChan = responses
	cs.RPCch <- rpc
	resp := <-responses
	if resp.Error != nil {
		return nil, resp.Error
	}
	rsp := new(AppendEntriesResponse)
	AppendEntriesResponseFromStruct(resp.Response.(*raft.AppendEntriesResponse), rsp)
	return rsp, nil
}
