package common

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func TimeFromProto(s *timestamp.Timestamp) time.Time {
	return s.AsTime()
}

func TimeToProto(s time.Time) *timestamp.Timestamp {
	return timestamppb.New(s)
}

func RPCHeaderFromProto(header *RPCHeader) raft.RPCHeader {
	h := raft.RPCHeader{}
	RPCHeaderToStruct(header, &h)
	return h
}

func RPCHeaderToProto(header raft.RPCHeader) *RPCHeader {
	h := new(RPCHeader)
	RPCHeaderToStruct(h, &header)
	return h
}

func LogFromProto(header []*Log) []*raft.Log {
	hs := make([]*raft.Log, len(header))
	for i, h := range hs {
		LogToStruct(header[i], h)
	}
	return hs
}

func LogToProto(header []*raft.Log) []*Log {
	hs := make([]*Log, len(header))
	for i, h := range hs {
		LogFromStruct(header[i], h)
	}
	return hs
}
