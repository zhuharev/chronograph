package chronograph

import (
	"fmt"
	"os"
	"testing"
)

var (
	testDBName = "test.bolt"
)

func clean() {
	os.RemoveAll(testDBName)
}

func TestTimeline(t *testing.T) {
	defer clean()
	var err error
	var bs Store
	bs, err = NewBoltStore(testDBName, true)
	if err != nil {
		t.Fatal(err.Error())
	}
	timeline := defaultTimeline{
		id:           "1",
		eventsIDSize: 8,
	}
	err = bs.TimelineCreate(timeline)
	if err != nil {
		t.Fatal(err.Error())
	}
	var events []Event
	var eventsLen = 10
	for i := 0; i < eventsLen; i++ {
		event := defaultEvent{
			id: fmt.Sprint(i),
		}
		err = bs.TimelineAppend(timeline, event)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	size, err := bs.TimelineSize(timeline)
	if err != nil {
		t.Fatal(err.Error())
	}

	if size != eventsLen {
		t.Fatalf("size not %d ,got %d", eventsLen, size)
	}

	// timeline

	events, err = bs.Timeline(timeline, 2, "7")
	if err != nil {
		t.Fatal(err.Error())
	}
	if l := len(events); l != 2 {
		t.Fatalf("len events unexpected need 2, got %d", l)
	}
	if events[0].ID() != "6" {
		t.Fatalf("id not 6, got %s", events[0].ID())
	}
}
