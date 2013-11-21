package raft

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestNetworkTransport_StartStop(t *testing.T) {
	trans, err := NewTCPTransport("127.0.0.1:0", time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	trans.Close()
}

func TestNetworkTransport_AppendEntries(t *testing.T) {
	// Transport 1 is consumer
	trans1, err := NewTCPTransport("127.0.0.1:0", time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer trans1.Close()
	rpcCh := trans1.Consumer()

	// Make the RPC request
	args := AppendEntriesRequest{
		Term:         10,
		Leader:       []byte("cartman"),
		PrevLogEntry: 100,
		PrevLogTerm:  4,
		Entries: []*Log{
			&Log{
				Index: 101,
				Term:  4,
				Type:  LogNoop,
			},
		},
		LeaderCommitIndex: 90,
	}
	resp := AppendEntriesResponse{
		Term:    4,
		LastLog: 90,
		Success: true,
	}

	// Listen for a request
	go func() {
		select {
		case rpc := <-rpcCh:
			// Verify the command
			req := rpc.Command.(*AppendEntriesRequest)
			if !reflect.DeepEqual(req, &args) {
				t.Fatalf("command mismatch: %#v %#v", *req, args)
			}

			rpc.Respond(&resp, nil)

		case <-time.After(200 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	// Transport 2 makes outbound request
	trans2, err := NewTCPTransport("127.0.0.1:0", time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer trans2.Close()

	var out AppendEntriesResponse
	if err := trans2.AppendEntries(trans1.LocalAddr(), &args, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the response
	if !reflect.DeepEqual(resp, out) {
		t.Fatalf("command mismatch: %#v %#v", resp, out)
	}
}

func TestNetworkTransport_RequestVote(t *testing.T) {
	// Transport 1 is consumer
	trans1, err := NewTCPTransport("127.0.0.1:0", time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer trans1.Close()
	rpcCh := trans1.Consumer()

	// Make the RPC request
	args := RequestVoteRequest{
		Term:         20,
		Candidate:    []byte("butters"),
		LastLogIndex: 100,
		LastLogTerm:  19,
	}
	resp := RequestVoteResponse{
		Term:    100,
		Peers:   []byte("blah"),
		Granted: false,
	}

	// Listen for a request
	go func() {
		select {
		case rpc := <-rpcCh:
			// Verify the command
			req := rpc.Command.(*RequestVoteRequest)
			if !reflect.DeepEqual(req, &args) {
				t.Fatalf("command mismatch: %#v %#v", *req, args)
			}

			rpc.Respond(&resp, nil)

		case <-time.After(200 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	// Transport 2 makes outbound request
	trans2, err := NewTCPTransport("127.0.0.1:0", time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer trans2.Close()

	var out RequestVoteResponse
	if err := trans2.RequestVote(trans1.LocalAddr(), &args, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the response
	if !reflect.DeepEqual(resp, out) {
		t.Fatalf("command mismatch: %#v %#v", resp, out)
	}
}

func TestNetworkTransport_InstallSnapshot(t *testing.T) {
	// Transport 1 is consumer
	trans1, err := NewTCPTransport("127.0.0.1:0", time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer trans1.Close()
	rpcCh := trans1.Consumer()

	// Make the RPC request
	args := InstallSnapshotRequest{
		Term:         10,
		Leader:       []byte("kyle"),
		LastLogIndex: 100,
		LastLogTerm:  9,
		Peers:        []byte("blah blah"),
		Size:         10,
	}
	resp := InstallSnapshotResponse{
		Term:    10,
		Success: true,
	}

	// Listen for a request
	go func() {
		select {
		case rpc := <-rpcCh:
			// Verify the command
			req := rpc.Command.(*InstallSnapshotRequest)
			if !reflect.DeepEqual(req, &args) {
				t.Fatalf("command mismatch: %#v %#v", *req, args)
			}

			// Try to read the bytes
			buf := make([]byte, 10)
			rpc.Reader.Read(buf)

			// Compare
			if bytes.Compare(buf, []byte("0123456789")) != 0 {
				t.Fatalf("bad buf %v", buf)
			}

			rpc.Respond(&resp, nil)

		case <-time.After(200 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	// Transport 2 makes outbound request
	trans2, err := NewTCPTransport("127.0.0.1:0", time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer trans2.Close()

	// Create a buffer
	buf := bytes.NewBuffer([]byte("0123456789"))

	var out InstallSnapshotResponse
	if err := trans2.InstallSnapshot(trans1.LocalAddr(), &args, &out, buf); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the response
	if !reflect.DeepEqual(resp, out) {
		t.Fatalf("command mismatch: %#v %#v", resp, out)
	}
}

func TestNetworkTransport_EncodeDecode(t *testing.T) {
	// Transport 1 is consumer
	trans1, err := NewTCPTransport("127.0.0.1:0", time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer trans1.Close()

	local := trans1.LocalAddr()
	enc := trans1.EncodePeer(local)
	dec := trans1.DecodePeer(enc)

	if dec.String() != local.String() {
		t.Fatalf("enc/dec fail: %v %v", dec, local)
	}
}
