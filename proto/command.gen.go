// Code generated by mog. DO NOT EDIT.

package commands

import "github.com/hashicorp/raft"

func AppendEntriesRequestToStruct(s *AppendEntriesRequest, t *raft.AppendEntriesRequest) {
	if s == nil {
		return
	}
	if s.RPCHeader != nil {
		RPCHeaderToStruct(s.RPCHeader, &t.RPCHeader)
	}
	t.Term = s.Term
	t.Leader = s.Leader
	t.PrevLogEntry = s.PrevLogEntry
	t.PrevLogTerm = s.PrevLogTerm
	{
		t.Entries = make([]*raft.Log, len(s.Entries))
		for i := range s.Entries {
			if s.Entries[i] != nil {
				var x raft.Log
				LogToStruct(s.Entries[i], &x)
				t.Entries[i] = &x
			}
		}
	}
	t.LeaderCommitIndex = s.LeaderCommitIndex
}
func AppendEntriesRequestFromStruct(t *raft.AppendEntriesRequest, s *AppendEntriesRequest) {
	if s == nil {
		return
	}
	{
		var x RPCHeader
		RPCHeaderFromStruct(&t.RPCHeader, &x)
		s.RPCHeader = &x
	}
	s.Term = t.Term
	s.Leader = t.Leader
	s.PrevLogEntry = t.PrevLogEntry
	s.PrevLogTerm = t.PrevLogTerm
	{
		s.Entries = make([]*Log, len(t.Entries))
		for i := range t.Entries {
			if t.Entries[i] != nil {
				var x Log
				LogFromStruct(t.Entries[i], &x)
				s.Entries[i] = &x
			}
		}
	}
	s.LeaderCommitIndex = t.LeaderCommitIndex
}
func AppendEntriesResponseToStruct(s *AppendEntriesResponse, t *raft.AppendEntriesResponse) {
	if s == nil {
		return
	}
	if s.RPCHeader != nil {
		RPCHeaderToStruct(s.RPCHeader, &t.RPCHeader)
	}
	t.Term = s.Term
	t.LastLog = s.LastLog
	t.Success = s.Success
	t.NoRetryBackoff = s.NoRetryBackoff
}
func AppendEntriesResponseFromStruct(t *raft.AppendEntriesResponse, s *AppendEntriesResponse) {
	if s == nil {
		return
	}
	{
		var x RPCHeader
		RPCHeaderFromStruct(&t.RPCHeader, &x)
		s.RPCHeader = &x
	}
	s.Term = t.Term
	s.LastLog = t.LastLog
	s.Success = t.Success
	s.NoRetryBackoff = t.NoRetryBackoff
}
func LogToStruct(s *Log, t *raft.Log) {
	if s == nil {
		return
	}
	t.Index = s.Index
	t.Term = s.Term
	t.Type = raft.LogType(s.Type)
	t.Data = s.Data
	t.Extensions = s.Extensions
	t.AppendedAt = TimeFromProto(s.AppendedAt)
}
func LogFromStruct(t *raft.Log, s *Log) {
	if s == nil {
		return
	}
	s.Index = t.Index
	s.Term = t.Term
	s.Type = uint32(t.Type)
	s.Data = t.Data
	s.Extensions = t.Extensions
	s.AppendedAt = TimeToProto(t.AppendedAt)
}
func RPCHeaderToStruct(s *RPCHeader, t *raft.RPCHeader) {
	if s == nil {
		return
	}
	t.ProtocolVersion = raft.ProtocolVersion(s.ProtocolVersion)
	t.ID = s.ID
	t.Addr = s.Addr
}
func RPCHeaderFromStruct(t *raft.RPCHeader, s *RPCHeader) {
	if s == nil {
		return
	}
	s.ProtocolVersion = int32(t.ProtocolVersion)
	s.ID = t.ID
	s.Addr = t.Addr
}
