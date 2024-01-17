// Package provides counter that receives path to resource with substring
// and counts non-overlapping instances of substr in resource.
package counter

import (
	"sync"
)

type counter struct {
	k       int
	total   int64
	wg      sync.WaitGroup
	limiter chan struct{}
}

// Inits the counter
func New(k int) Counter {
	return &counter{
		k,
		0,
		sync.WaitGroup{},
		make(chan struct{}, k),
	}
}
