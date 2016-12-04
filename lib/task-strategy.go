package lib

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/types"
)

// serverGroup dynamic hold the ranked priorities of all smtp server
// according by the smtp delivery results, it provide the ability for
// recovering from temporary un-available servers.
type serverGroup struct {
	raw   []*types.SMTPServer
	mux   sync.RWMutex // protect pools
	pools map[bson.ObjectId]*rankedServer
}

type rankedServer struct {
	rank     uint64 // weight of current server
	server   *types.SMTPServer
	positive uint64
	negative uint64
}

func newServerGroup(ss []*types.SMTPServer) *serverGroup {
	sg := &serverGroup{
		raw:   ss,
		pools: make(map[bson.ObjectId]*rankedServer),
	}
	for _, s := range ss {
		sg.pools[s.ID] = &rankedServer{1, s, 0, 0}
	}
	return sg
}

func (sg serverGroup) stats() map[string]interface{} {
	sg.mux.RLock()
	defer sg.mux.RUnlock()
	var (
		m = map[string]interface{}{}
	)
	for id, rs := range sg.pools {
		m[id.Hex()] = map[string]interface{}{
			"rank":     rs.rank,
			"positive": rs.positive,
			"negative": rs.negative,
		}
	}
	return m
}

func (sg serverGroup) next() *types.SMTPServer {
	sg.mux.RLock()
	defer sg.mux.RUnlock()
	rand.Seed(time.Now().UnixNano())
	var (
		r = uint64(rand.Int63n(int64(sg.rankSum())))
		n = uint64(0)
	)
	for _, rs := range sg.pools {
		n += rs.rank
		if n > r {
			return rs.server
		}
	}
	return nil
}

func (sg serverGroup) rankSum() uint64 {
	var sum uint64
	for _, rs := range sg.pools {
		sum += rs.rank
	}
	return sum
}

func (sg serverGroup) upgrade(id bson.ObjectId) {
	sg.mux.Lock()
	if v, ok := sg.pools[id]; ok {
		sg.pools[id].positive++
		if v.rank < math.MaxUint64 { // rank maximize is math.MaxUint64
			sg.pools[id].rank++
		}
	}
	sg.mux.Unlock()
}

func (sg serverGroup) downgrade(id bson.ObjectId) {
	sg.mux.Lock()
	if v, ok := sg.pools[id]; ok {
		sg.pools[id].negative++
		if v.rank > 1 { // rank minimal is 1
			sg.pools[id].rank--
		}
	}
	sg.mux.Unlock()
}

func (sg serverGroup) size() int {
	return len(sg.raw)
}

func (sg serverGroup) empty() bool {
	return sg.size() == 0
}
