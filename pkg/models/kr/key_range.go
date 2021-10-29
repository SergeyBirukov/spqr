package kr

import (
	"github.com/pg-sharding/spqr/qdb/qdb"
	proto "github.com/pg-sharding/spqr/router/protos"
)

type KeyRangeBound []byte

type ShardKey struct {
	Name string
	RW   bool
}

type KeyRange struct {
	LowerBound []byte
	UpperBound []byte
	Shid       string
	ID         string
}

func CmpRanges(kr []byte, other []byte) bool {
	if len(kr) == len(other) {
		return string(kr) <= string(other)
	}

	return len(kr) <= len(other)
}

func KeyRangeFromSQL(kr *qdb.KeyRange) *KeyRange {
	return &KeyRange{
		LowerBound: kr.From,
		UpperBound: kr.To,
		Shid:       kr.ShardID,
		ID:         kr.KeyRangeID,
	}
}

func (kr *KeyRange) ToSQL() *qdb.KeyRange {
	return &qdb.KeyRange{
		From:       kr.LowerBound,
		To:         kr.UpperBound,
		ShardID:    kr.Shid,
		KeyRangeID: kr.ID,
	}
}

func (kr *KeyRange) ToProto() *proto.KeyRange {
	return &proto.KeyRange{
		LowerBound: string(kr.LowerBound),
		UpperBound: string(kr.UpperBound),
		ShardId:    kr.Shid,
		Krid:       kr.ID,
	}
}

func KeyRangeFromProto(kr *proto.KeyRange) *KeyRange {
	if kr == nil {
		return nil
	}
	return &KeyRange{
		LowerBound: []byte(kr.LowerBound),
		UpperBound: []byte(kr.UpperBound),
		Shid:       kr.ShardId,
		ID:         kr.Krid,
	}
}
