// Package provides counter that receives path to resource with substring
// and counts non-overlapping instances of substr in resource.
package counter

import (
	"net/http"
	"sync"
	"time"
)

type counter struct {
	k       int
	jobs    int
	total   int64
	wg      sync.WaitGroup
	limiter chan struct{}
	counts  chan counted
	client  http.Client
}

// Inits the counter
func New(k int) Counter {
	return &counter{
		k,
		0,
		0,
		sync.WaitGroup{},
		make(chan struct{}, k),
		make(chan counted),
		http.Client{
			Timeout: 5 * time.Second,
		},
	}
}
