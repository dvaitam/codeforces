package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	letters := []rune{'A', 'B', 'C'}
	for ; t > 0; t-- {
		grid := make([]string, 3)
		for i := 0; i < 3; i++ {
			fmt.Fscan(in, &grid[i])
		}
		for i := 0; i < 3; i++ {
			if strings.ContainsRune(grid[i], '?') {
				seen := map[rune]bool{}
				for _, ch := range grid[i] {
					if ch != '?' {
						seen[ch] = true
					}
				}
				for _, l := range letters {
					if !seen[l] {
						fmt.Fprintln(out, string(l))
						break
					}
				}
				break
			}
		}
	}
}
