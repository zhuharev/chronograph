package chronograph

import (
	"bytes"
	"sync"

	"github.com/Unknwon/com"
	"github.com/zhuharev/boltutils"
	"github.com/zhuharev/intarr"
)

type Store interface {
	TimelineCreate(Timeline) error
	Timeline(timeline Timeline, limit int, startIDs ...string) ([]Event, error)
	TimelineAppend(Timeline, Event) error
	TimelineSize(Timeline) (int, error)
}

var (
	stores = map[string]NewStoreFunc{}
	mx     = sync.Mutex{}
)

func RegistreSore(name string, fn NewStoreFunc) {
	mx.Lock()
	stores[name] = fn
	mx.Unlock()
}

type NewStoreFunc func(string) (Store, error)

func NewStore(storeName, storeOptions string) (Store, error) {
	return stores[storeName](storeOptions)
}

var (
	boltPrefix                   = "chronograph_"
	boltTimelineBucketName       = []byte(boltPrefix + "timelines")
	boltTimelineEventsBucketName = []byte("events")
)

// buckets logic: /timelines/:id/events
//                /timelines/:id/events_data etc
type boltStore struct {
	db                 *boltutils.DB
	compressionEnabled bool
}

func init() {
	RegistreSore("bolt", func(setting string) (Store, error) {
		return NewBoltStore(setting, true)
	})
}

func NewBoltStore(path string, enableCompression bool) (bs *boltStore, err error) {
	db, err := boltutils.Open(path, 0777, nil)
	return &boltStore{db: db, compressionEnabled: enableCompression}, err
}

func (bs *boltStore) TimelineCreate(t Timeline) (err error) {
	err = bs.db.CreateBucketPath(makeTimelineEventsPath(t))
	if err != nil {
		return err
	}
	//err = bs.db.CreateBucketPath(makeTimelineEventsPathWithKey(t))
	return err
}

// Timeline get events by timeline
// in boltdb all keys saved as single bytes slice and all data
// will be recieved in each call this func
// if startIDs has an element, events will be returned after ithes ids with
// limit
func (bs *boltStore) Timeline(timeline Timeline, limit int, startIDs ...string) (events []Event, err error) {
	var (
		data []byte
	)
	if bs.compressionEnabled {
		data, err = bs.db.GetGzipped(makeTimelineEventsPathWithKey(timeline))
	} else {
		data, err = bs.db.GetPath(makeTimelineEventsPathWithKey(timeline))
	}
	if err != nil {
		if err != boltutils.ErrNotFound {
			return
		}
		err = nil
	}

	// TODO: use bufio
	var (
		l   = len(data) / timeline.EventsIDSize()
		buf = bytes.NewReader(data)
		tmp = make([]byte, timeline.EventsIDSize())

		needSkip = len(startIDs) > 0
	)
	for i := 0; i < l; i++ {
		if timeline.OrderChronologic() {
			_, err = buf.Read(tmp)
		} else {
			var offset = int64(timeline.EventsIDSize() * (l - i - 1))
			_, err = buf.ReadAt(tmp, offset)
		}
		if err != nil {
			return
		}
		currentEqual := bytes.EqualFold([]byte(startIDs[0]), trimBytes(tmp))
		if needSkip && !currentEqual {
			continue
		}
		needSkip = false
		// clear empty head

		if !currentEqual {
			events = append(events, defaultEvent{string(trimBytes(tmp)), false})
		}
		if len(events) == limit {
			break
		}
	}

	return
}

func (bs *boltStore) TimelineAppend(timeline Timeline, ev Event) (err error) {
	var (
		data []byte
		path = makeTimelineEventsPath(timeline)
	)

	// convert event id to sized slice
	// make empty slice
	// check if id intable(numeric)
	bytesID := []byte(ev.ID())
	if ev.IDIntable() {
		bytesID = intarr.Uint64ToBytes(uint64(com.StrTo(ev.ID()).MustInt64()))
	}

	idLen := len(bytesID)
	emptyHead := make([]byte, timeline.EventsIDSize()-idLen)
	idKey := append(emptyHead, bytesID...)

	if bs.compressionEnabled {
		data, err = bs.db.GetGzipped(path, boltTimelineEventsBucketName)
		if err != nil {
			if err != boltutils.ErrNotFound {
				return err
			}
		}
		return bs.db.PutGzip(path, boltTimelineEventsBucketName, append(data, idKey...))
	}
	data, err = bs.db.GetPath(path, boltTimelineEventsBucketName)
	if err != nil {
		if err != boltutils.ErrNotFound {
			return err
		}
	}
	return bs.db.Put(path, boltTimelineEventsBucketName, append(data, idKey...))
}

func (bs *boltStore) TimelineSize(timeline Timeline) (int, error) {
	var (
		data []byte
		err  error
	)
	if bs.compressionEnabled {
		data, err = bs.db.GetGzipped(makeTimelineEventsPathWithKey(timeline))
	} else {
		data, err = bs.db.GetPath(makeTimelineEventsPathWithKey(timeline))
	}
	if err != nil {
		return 0, err
	}
	return len(data) / timeline.EventsIDSize(), nil
}
