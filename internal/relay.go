package internal

import (
	"github.com/jackc/pgproto3"
	"github.com/pg-sharding/spqr/internal/qrouter"
	"github.com/wal-g/tracelog"
)

type RelayState struct {
	TxActive bool

	ActiveShards []ShardKey
}

func (rst *RelayState) reroute(rt qrouter.Qrouter, cl Client, cmngr ConnManager, q *pgproto3.Query) error {

	shardKeys := rt.Route(q.String)

	var shards []ShardKey

	for _, name := range shardKeys {
		shards = append(shards, NewSHKey(name))
	}

	if shardKeys == nil {
		return cmngr.UnRouteWithError(cl, nil, "failed to match shard")
	}
	tracelog.InfoLogger.Printf("parsed shard names %s", shardKeys)

	if err := cmngr.UnRouteCB(cl, rst.ActiveShards); err != nil {
		tracelog.ErrorLogger.PrintError(err)
	}
	rst.ActiveShards = shards

	var serv Server
	var err error

	if len(shards) > 1 {
		serv, err = NewMultiShardServer(cl.Route().beRule, cl.Route().servPool, shards)
		if err != nil {
			return err
		}
	} else {
		serv = NewShardServer(cl.Route().beRule, cl.Route().servPool)
	}

	_ = cl.AssignServerConn(serv)

	tracelog.InfoLogger.Printf("route cl %s:%s to %v", cl.Usr(), cl.DB(), shards)

	if err := cmngr.RouteCB(cl, rst.ActiveShards); err != nil {
		tracelog.ErrorLogger.PrintError(err)
	}

	return nil
}

const TXREL = 73

func (rst *RelayState) completeRelay(cl Client, cmngr ConnManager, txst byte) error {

	if err := cl.Send(&pgproto3.ReadyForQuery{}); err != nil {
		return err
	}

	if txst == TXREL {
		if rst.TxActive {
			if err := cmngr.TXEndCB(cl, rst); err != nil {
				return err
			}
			rst.TxActive = false
		}
	} else {
		if !rst.TxActive {
			if err := cmngr.TXBeginCB(cl, rst); err != nil {
				return err
			}
			rst.TxActive = true
		}
	}

	return nil
}
