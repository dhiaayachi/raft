package commands

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func TimeFromProto(s *timestamp.Timestamp) time.Time {
	return s.AsTime()
}

func TimeToProto(s time.Time) *timestamp.Timestamp {
	return timestamppb.New(s)
}
