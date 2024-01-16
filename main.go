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
	)

	for scanner.Scan() {
		counter.Count(scanner.Text(), *substr)
	}

	fmt.Printf("Total: %d\n", counter.GetTotal())
}
