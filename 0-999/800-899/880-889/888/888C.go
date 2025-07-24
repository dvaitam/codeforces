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

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	n := len(s)
	best := n + 1
	for ch := byte('a'); ch <= byte('z'); ch++ {
		last := -1
		maxDiff := 0
		for i := 0; i < n; i++ {
			if s[i] == ch {
				diff := i - last
				if diff > maxDiff {
					maxDiff = diff
				}
				last = i
			}
		}
		diff := n - last
		if diff > maxDiff {
			maxDiff = diff
		}
		if maxDiff < best {
			best = maxDiff
		}
	}

	fmt.Fprintln(writer, best)
}
