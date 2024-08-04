package snap

import (
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/eth/protocols/snap"
	"github.com/zhiqiangxu/litenode"
	"github.com/zhiqiangxu/litenode/eth/common"
)

type Syncer struct {
	sync.Mutex
	path     string
	statusCh <-chan eth.MinStatus
	doneCh   chan struct{}
	wg       sync.WaitGroup
}

func NewSyncer(path string) *Syncer {
	return &Syncer{path: path, doneCh: make(chan struct{})}
}

func (s *Syncer) Register(peer *snap.Peer) error {
	s.Lock()
	defer s.Unlock()

	return nil
}

func (s *Syncer) Unregister(peer *snap.Peer) error {
	s.Lock()
	defer s.Unlock()

	return nil
}

func (s *Syncer) Start(node *litenode.Node, statusCh <-chan eth.MinStatus, snapMsgCh chan common.SnapSyncPacket) {
	s.statusCh = statusCh
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		var tenBiggestStatus biggestStatus
		var status eth.MinStatus
		for {
			select {
			case status = <-s.statusCh:
				tenBiggestStatus.Push(status, 10)
			case <-s.doneCh:
				return
			}
		}
	}()
}

func (s *Syncer) Stop() {
	close(s.doneCh)
	s.wg.Wait()
}

type biggestStatus []*eth.MinStatus

func (b *biggestStatus) Push(status eth.MinStatus, k int) {
	ix := sort.Search(len(*b), func(i int) bool {
		return status.TD.Cmp((*b)[i].TD) > 0
	})
	if ix == len(*b) {
		// farther away than all nodes we already have.
		// if there was room for it, the node is now the last element.
		if len(*b) < k {
			*b = append(*b, &status)
		}
	} else {
		// slide existing entries down to make room
		// this will overwrite the entry we just appended.
		var last *eth.MinStatus
		if len(*b) < k {
			last = (*b)[len(*b)-1]
		}
		copy((*b)[ix+1:], (*b)[ix:])
		(*b)[ix] = &status
		if last != nil {
			*b = append(*b, last)
		}
	}

}
