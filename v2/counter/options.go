package counter

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type counted struct {
	path  string
	count int
}

// checks if path is url or file
func isURL(path string) bool {
	return strings.Contains(path, "http")
}

// Increases total count
func (c *counter) incTotal(delta int) {
	atomic.AddInt64(&c.total, int64(delta))
}

// Counts instances of substring in given reader
func (c *counter) count(path, substr string, data io.Reader) int {
	var (
		count   int
		scanner = bufio.NewScanner(data)
	)
	for scanner.Scan() {
		count += strings.Count(scanner.Text(), substr)
	}
	return count
}

// Starts counting worker
func (c *counter) startCountWorker(path, substr string) (int, error) {
	if isURL(path) {
		resp, err := http.Get(path)
		if err != nil {
			return 0, fmt.Errorf("failed to get data from %s (%v)", path, err)
		}
		defer resp.Body.Close()
		return c.count(path, substr, resp.Body), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s (%v)", path, err)
	}
	defer file.Close()
	return c.count(path, substr, file), nil
}

// Adds a task to the counter
func (c *counter) Count(path, substr string) {
	c.wg.Add(1)
	c.jobs++
	go func() {
		c.limiter <- struct{}{}
		defer func() {
			c.wg.Done()
			<-c.limiter
		}()
		c.startCountWorker(path, substr)
		count, err := c.startCountWorker(path, substr)
		defer func() { go func() { c.counts <- counted{path, count} }() }()
		if err != nil {
			slog.Error(err.Error())
			return
		}
		c.incTotal(count)
	}()
}

// Gets total count and resets counter
func (c *counter) GetTotal() int64 {
	defer func() {
		c.total = 0
		c.jobs = 0
		for {
			select {
			case <-c.counts:
				time.Sleep(time.Millisecond * 10)
			default:
				return
			}
		}
	}()
	c.wg.Wait()
	return c.total
}

// Gets {path, count} in order of channel receiving
func (c *counter) GetCount() (string, int) {
	if c.jobs == 0 {
		return "", 0
	}
	c.jobs--
	counted := <-c.counts
	return counted.path, counted.count
}
