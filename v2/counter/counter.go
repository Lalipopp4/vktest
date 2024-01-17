package counter

// Exported counter interaface
type Counter interface {

	// Adds a task to the counter
	Count(path, substr string)

	// Gets total count and resets counter
	GetTotal() int64

	// Gets {path, count} in order of channel receiving (receives empty if called more than got counts)
	GetCount() (string, int)
}
