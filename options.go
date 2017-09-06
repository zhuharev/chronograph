package chronograph

type Options struct {
	StoreName    string
	StoreOptions string
}

func getOptions(opts ...Options) Options {
	if len(opts) > 0 {
		return opts[0]
	}
	return Options{
		StoreName:    "bolt",
		StoreOptions: "chronograph.bolt",
	}
}
