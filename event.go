package chronograph

type Event interface {
	IDer
}

type Eventer interface {
	ToEvent() Event
}

type defaultEvent struct {
	id        string
	idIntable bool
}

func NewEvent(id string, idIntable bool) Event {
	return defaultEvent{id, idIntable}
}

func (de defaultEvent) ID() string {
	return de.id
}

func (de defaultEvent) IDIntable() bool {
	return de.idIntable
}
