package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	maxT := 2 * n
	allowed := make([]bool, maxT+1)
	for i := 0; i < k; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		for t := l; t <= r && t <= maxT; t++ {
			allowed[t] = true
		}
	}

	inf := int(1e9)
	offset := n
	curA := make([]int, 2*n+1)
	curB := make([]int, 2*n+1)
	for i := range curA {
		curA[i] = inf
		curB[i] = inf
	}
	curA[offset] = 0

	for t := 0; t < maxT; t++ {
		nextA := make([]int, 2*n+1)
		nextB := make([]int, 2*n+1)
		for i := range nextA {
			nextA[i] = inf
			nextB[i] = inf
		}
		for d := 0; d <= 2*n; d++ {
			if curA[d] < inf {
				if d+1 <= 2*n {
					nextA[d+1] = min(nextA[d+1], curA[d])
				}
				if allowed[t] && d-1 >= 0 {
					nextB[d-1] = min(nextB[d-1], curA[d]+1)
				}
			}
			if curB[d] < inf {
				if d-1 >= 0 {
					nextB[d-1] = min(nextB[d-1], curB[d])
				}
				if allowed[t] && d+1 <= 2*n {
					nextA[d+1] = min(nextA[d+1], curB[d]+1)
				}
			}
		}
		curA = nextA
		curB = nextB
	}

	ans := min(curA[offset], curB[offset])
	if ans >= inf {
		fmt.Fprintln(out, "Hungry")
	} else {
		fmt.Fprintln(out, "Full")
		fmt.Fprintln(out, ans)
	}
}
