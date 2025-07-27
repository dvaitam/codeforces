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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// nextIdx[i] is next position > i where a[i] occurs, or large value if none
	nextIdx := make([]int, n)
	lastPos := make(map[int]int)
	const inf = int(1e9)
	for i := n - 1; i >= 0; i-- {
		if p, ok := lastPos[a[i]]; ok {
			nextIdx[i] = p
		} else {
			nextIdx[i] = inf
		}
		lastPos[a[i]] = i
	}

	lastA, lastB := -1, -1
	posA, posB := inf, inf
	segA, segB := 0, 0

	for i := 0; i < n; i++ {
		x := a[i]
		nxt := nextIdx[i]
		if x == lastA {
			posA = nxt
			continue
		}
		if x == lastB {
			posB = nxt
			continue
		}
		// need to place x in one of the sequences
		if posA <= posB { // keep A, replace B
			if x != lastB {
				segB++
				lastB = x
			}
			posB = nxt
		} else { // replace A
			if x != lastA {
				segA++
				lastA = x
			}
			posA = nxt
		}
	}

	fmt.Fprintln(out, segA+segB)
}
