package trace

type options struct {
	env string
}

type Option interface {
	apply(*options)
}

type envOption string

func (e envOption) apply(opts *options) {
	opts.env = string(e)
}

// WithEnvironment configures the environment to tag traces/spans with.
func WithEnvironment(env string) Option {
	return envOption(env)
}
