package chronograph

type IDer interface {
	ID() string
	IDIntable() bool
}

type Timeline interface {
	IDer
	OrderChronologic() bool // chronologic or antichronologic todo: enum

	// metadata
	EventsIDSize() int // fixrd size for timeline events

	// EventsHasData() bool
}

type defaultTimeline struct {
	id                string
	orderCronological bool
	eventsIDSize      int
	eventsIDIntable   bool
}

func (dt defaultTimeline) ID() string {
	return dt.id
}

func (dt defaultTimeline) IDIntable() bool {
	return true
}

func (dt defaultTimeline) OrderChronologic() bool {
	return dt.orderCronological
}

func (dt defaultTimeline) EventsIDSize() int {
	return dt.eventsIDSize
}

type Event interface {
	IDer
}

type defaultEvent struct {
	id string
}

func (de defaultEvent) ID() string {
	return de.id
}

func (de defaultEvent) IDIntable() bool {
	return false
}
