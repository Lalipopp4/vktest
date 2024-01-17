package counter

// Exported counter interaface
type Counter interface {

	// Adds a task to the counter
	Count(path, substr string)

	// Gets total count and resets counter
	GetTotal() int64

	// Gets {path, count} in order of channel receiving (can block goroutine if called too many times)
	GetCount() (string, int)
}
