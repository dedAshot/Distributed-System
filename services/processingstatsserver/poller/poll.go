package poller

import (
	"container/list"
	"fmt"
	"sort"
	"statisticsserver/store"
	"sync"
	"time"
)

var (
	stopPolling         = false
	defaultPollingCount = 10000
	SavedStats          = make([]*store.StatRow, 0, 10000)
)

type QueryIndex struct {
	startId int64
	count   int64
}

func NewQueryIndex(startId, count int64) *QueryIndex {
	return &QueryIndex{startId: startId, count: count}
}

func StartPolling(interval time.Duration) {
	go poll(interval)
}

func StopPolling() {
	stopPolling = true
}

func poll(interval time.Duration) {
	for !stopPolling {
		pollDb()
		time.Sleep(interval)
	}
}

func pollDb() {

	stats, err := store.GetLastMessages(defaultPollingCount)
	if err != nil {
		fmt.Println(err)
	}

	SavedStats = stats
}

func GetPage(startId, rowCount int) ([]*store.StatRow, error) {

	if startId < 0 {
		l := len(SavedStats)
		if l < rowCount {
			return SavedStats[:l], nil
		} else {
			return SavedStats[:rowCount], nil
		}
	}

	chacheIndex := sort.Search(len(SavedStats), func(id int) bool {
		if SavedStats[id].Id <= startId { // <= SavedStats is in desc
			return true
		} else {
			return false
		}
	})

	if chacheIndex < len(SavedStats) && SavedStats[chacheIndex].Id == startId {
		if chacheIndex-rowCount < 0 {
			return SavedStats[:chacheIndex+1], nil
		} else {
			return SavedStats[chacheIndex-rowCount+1:chacheIndex+1], nil
		}

	} else {
		return store.GetMessages(startId, rowCount)
	}
}

type chacheEl struct {
	stat []*store.StatRow
	key  QueryIndex
}

func NewChacheEl(stat []*store.StatRow, key QueryIndex) *chacheEl {
	return &chacheEl{stat: stat, key: key}
}

var (
	DbStatsChaches = make(map[QueryIndex][]*store.StatRow)
)

type DbStatsChache struct {
	chache     map[QueryIndex]*list.Element
	maxElCount int
	lru        list.List
	rwm        sync.RWMutex
}

func (ch *DbStatsChache) StoreInChache(stats []*store.StatRow, key QueryIndex) error {
	chEl := NewChacheEl(stats, key)

	ch.rwm.Lock()
	if ch.lru.Len() >= ch.maxElCount {
		lastEl := ch.lru.Remove(ch.lru.Back()).(*chacheEl)
		lkey := lastEl.key
		delete(ch.chache, lkey)
	}

	ch.chache[key] = ch.lru.PushFront(chEl)

	return nil
}

func (ch *DbStatsChache) GetChache(startId, count int64) ([]*store.StatRow, bool) {
	ch.rwm.RLock()

	elem, ok := ch.chache[*NewQueryIndex(startId, count)]
	if !ok {
		return nil, false
	}

	ch.rwm.RUnlock()

	ch.lru.MoveToFront(elem)
	stats := elem.Value.([]*store.StatRow)

	return stats, true
}
