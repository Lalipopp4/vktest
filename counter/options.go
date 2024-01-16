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
)

// checks if path is url or file
func isURL(path string) bool {
	return strings.Contains(path, "http")
}

// Increases total count
func (c *counter) incTotal(delta int) {
	atomic.AddInt64(&c.total, int64(delta))
}

// Counts instances of substring in given reader
func (c *counter) count(source, substr string, data io.Reader) int {
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
	var count int

	if isURL(path) {
		resp, err := http.Get(path)
		if err != nil {
			return 0, fmt.Errorf("failed to get data from %s (%v)", path, err)
		}
		defer resp.Body.Close()
		count = c.count(path, substr, resp.Body)
		c.incTotal(count)
		return count, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s (%v)", path, err)
	}
	defer file.Close()
	count = c.count(path, substr, file)
	c.incTotal(count)
	return count, nil
}

// Adds a task to the counter
func (c *counter) Count(path, substr string) {
	c.wg.Add(1)
	go func() {
		c.limiter <- struct{}{}
		defer c.wg.Done()
		defer func() { <-c.limiter }()
		count, err := c.startCountWorker(path, substr)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		fmt.Printf("Count for %s: %d\n", path, count)
	}()
}

// Gets total count and resets counter
func (c *counter) GetTotal() int64 {
	defer func() { c.total = 0 }()
	c.wg.Wait()
	return c.total
}
