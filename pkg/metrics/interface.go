package metrics

import (
	"github.com/uber-go/tally"
)

var DefaultTallyBuckets = tally.ValueBuckets{.01, .05, .1, .2, .5, .8, .9, 1, 5, 10, 15, 30, 60, 90, 120}
