package dual_transport

import (
	"github.com/hashicorp/raft"
	"io"
)

type DualTransport struct {
	addr      raft.ServerAddress
	consumer  chan raft.RPC
	primary   raft.Transport
	secondary raft.Transport
}

func (g DualTransport) Consumer() <-chan raft.RPC {
	return g.consumer
}

func (g DualTransport) LocalAddr() raft.ServerAddress {
	return g.addr
}

func (g DualTransport) AppendEntriesPipeline(id raft.ServerID, target raft.ServerAddress) (raft.AppendPipeline, error) {
	ap, err := g.primary.AppendEntriesPipeline(id, target)
	if err != nil {
		return g.secondary.AppendEntriesPipeline(id, target)
	}
	return ap, nil
}

func (g DualTransport) AppendEntries(id raft.ServerID, target raft.ServerAddress, args *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) error {
	if err := g.primary.AppendEntries(id, target, args, resp); err != nil {
		return g.secondary.AppendEntries(id, target, args, resp)
	}
	return nil
}

func (g DualTransport) RequestVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestVoteRequest, resp *raft.RequestVoteResponse) error {
	if err := g.primary.RequestVote(id, target, args, resp); err != nil {
		return g.secondary.RequestVote(id, target, args, resp)
	}
	return nil
}

func (g DualTransport) InstallSnapshot(id raft.ServerID, target raft.ServerAddress, args *raft.InstallSnapshotRequest, resp *raft.InstallSnapshotResponse, data io.Reader) error {
	if err := g.primary.InstallSnapshot(id, target, args, resp, data); err != nil {
		return g.secondary.InstallSnapshot(id, target, args, resp, data)
	}
	return nil
}

func (g DualTransport) EncodePeer(id raft.ServerID, addr raft.ServerAddress) []byte {
	return g.primary.EncodePeer(id, addr)
}

func (g DualTransport) DecodePeer(bytes []byte) raft.ServerAddress {
	return g.primary.DecodePeer(bytes)
}

func (g DualTransport) SetHeartbeatHandler(cb func(rpc raft.RPC)) {
	g.primary.SetHeartbeatHandler(cb)
	g.secondary.SetHeartbeatHandler(cb)
}

func (g DualTransport) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	if err := g.primary.TimeoutNow(id, target, args, resp); err != nil {
		return g.secondary.TimeoutNow(id, target, args, resp)
	}
	return nil
}

// NewDualTransport is used to initialize a new transport and
// generates a random local address if none is specified. The given timeout
// will be used to decide how long to wait for a connected peer to process the
// RPCs that we're sending it. See also Connect() and Consumer().
func NewDualTransport(primary raft.Transport, secondary raft.Transport) (raft.ServerAddress, raft.Transport) {
	pc := primary.Consumer()
	sc := secondary.Consumer()
	rpcs := make(chan raft.RPC)
	go func() {
		for {
			rpc := <-pc
			rpcs <- rpc
		}
	}()

	go func() {
		for {
			rpc := <-sc
			rpcs <- rpc
		}
	}()

	addr := primary.LocalAddr() + "%" + secondary.LocalAddr()
	return addr, &DualTransport{primary: primary, addr: addr, consumer: rpcs, secondary: secondary}

}
