package chronograph

type Chronograph struct {
	store Store
}

func New(opts ...Options) (chrono *Chronograph, err error) {
	chrono = &Chronograph{}
	opt := getOptions(opts...)
	chrono.store, err = NewStore(opt.StoreName, opt.StoreOptions)
	return
}

func (chrono *Chronograph) CreateTimeline(timeliner Timeliner) error {
	return chrono.store.TimelineCreate(timeliner.ToTimeline())
}

func (chrono *Chronograph) Append(timeliner Timeliner, eventer Eventer) (err error) {
	err = chrono.store.TimelineAppend(timeliner.ToTimeline(), eventer.ToEvent())
	return
}
