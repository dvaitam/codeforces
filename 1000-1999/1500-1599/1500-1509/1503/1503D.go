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

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i], &b[i])
	}

	if n > 16 {
		fmt.Fprintln(out, -1)
		return
	}

	m := 1 << n
	minFlips := -1
	for mask := 0; mask < m; mask++ {
		front := make([]int, n)
		back := make([]int, n)
		flips := 0
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				front[i] = b[i]
				back[i] = a[i]
				flips++
			} else {
				front[i] = a[i]
				back[i] = b[i]
			}
		}
		// sort indices by front values
		idx := make([]int, n)
		for i := 0; i < n; i++ {
			idx[i] = i
		}
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if front[idx[i]] > front[idx[j]] {
					idx[i], idx[j] = idx[j], idx[i]
				}
			}
		}
		ok := true
		for i := 1; i < n; i++ {
			if back[idx[i-1]] <= back[idx[i]] {
				ok = false
				break
			}
		}
		if ok {
			if minFlips == -1 || flips < minFlips {
				minFlips = flips
			}
		}
	}
	fmt.Fprintln(out, minFlips)
}
