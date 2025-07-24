package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var s string
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	freq := make(map[string]int)
	best := ""
	bestCount := 0
	for i := 0; i < n-1; i++ {
		two := s[i : i+2]
		freq[two]++
		if freq[two] > bestCount {
			bestCount = freq[two]
			best = two
		}
	}

	fmt.Fprintln(writer, best)
}
