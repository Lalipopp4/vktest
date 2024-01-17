package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"vktest/counter"
)

func main() {
	k := flag.Int("k", 5, "Maximum number of workers")
	substr := flag.String("substring", "go", "Substring to search")
	flag.Parse()

	var (
		counter = counter.New(*k)
		scanner = bufio.NewScanner(os.Stdin)
		count   int
	)

	for scanner.Scan() {
		counter.Count(scanner.Text(), *substr)
		count++
	}

	for i := 0; i < count; i++ {
		path, count := counter.GetCount()
		fmt.Printf("Count for %s: %d\n", path, count)
	}

	fmt.Printf("Total: %d\n", counter.GetTotal())
}
