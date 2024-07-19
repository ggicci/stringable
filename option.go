package stringable

type Option func(o *options)

func NoHybrid() Option {
	return func(o *options) {
		o.Opt(optionNoHybrid)
	}
}

func CompleteHybrid() Option {
	return func(o *options) {
		o.Opt(optionCompleteHybrid)
	}
}

type options struct {
	Value uint8
}

func defaultOptions() *options {
	return &options{}
}

func (o *options) Opt(v option) {
	o.Value |= uint8(v)
}

func (o *options) Has(v option) bool {
	return (o.Value & uint8(v)) > 0
}

type option int

const (
	optionNoHybrid option = 1 << iota
	optionCompleteHybrid
)
