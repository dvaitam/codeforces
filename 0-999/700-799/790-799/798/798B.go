package main

import (
	"bufio"
	"fmt"
	"os"
)

func rotate(s string, k int) string {
	n := len(s)
	k %= n
	return s[k:] + s[:k]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	strs := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &strs[i])
	}
	length := len(strs[0])
	const inf = int(1e9)
	best := inf

	for k := 0; k < length; k++ {
		target := rotate(strs[0], k)
		total := k
		possible := true
		for i := 1; i < n && possible; i++ {
			found := -1
			for j := 0; j < length; j++ {
				if rotate(strs[i], j) == target {
					found = j
					break
				}
			}
			if found == -1 {
				possible = false
				break
			}
			total += found
		}
		if possible && total < best {
			best = total
		}
	}

	if best == inf {
		fmt.Println(-1)
	} else {
		fmt.Println(best)
	}
}
