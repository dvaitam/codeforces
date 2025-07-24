package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	present := make(map[int]bool)
	cost := 0

	for i := 0; i < n; i++ {
		book := a[i]
		if !present[book] {
			cost++
			if len(present) == k {
				// choose a book to remove
				removeBook := 0
				farthest := -1
				for b := range present {
					next := n + 1
					for j := i + 1; j < n; j++ {
						if a[j] == b {
							next = j
							break
						}
					}
					if next > farthest {
						farthest = next
						removeBook = b
					}
				}
				delete(present, removeBook)
			}
			present[book] = true
		}
	}
	fmt.Fprintln(out, cost)
}
