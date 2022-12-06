package publisher

type Options struct {
	continueOnErr bool
	countFailure  bool
	logFailure    bool
	failure       int
}

// NewOptions return an option struct for publisher
func NewOptions(continueOnErr, countFailure, logFailure bool) *Options {
	return &Options{
		continueOnErr: continueOnErr,
		countFailure:  countFailure,
		logFailure:    logFailure,
		failure:       0,
	}
}
