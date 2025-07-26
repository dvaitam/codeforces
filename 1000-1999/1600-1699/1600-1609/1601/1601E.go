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

	var n, q, k int
	if _, err := fmt.Fscan(reader, &n, &q, &k); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		a[i] = x
	}

	// Precompute sliding window minima of length k for each day
	windowMin := make([]int64, n+1)
	deque := make([]int, 0)
	for i := 1; i <= n; i++ {
		// remove indices out of window
		if len(deque) > 0 && deque[0] <= i-k {
			deque = deque[1:]
		}
		// maintain increasing order in deque
		for len(deque) > 0 && a[deque[len(deque)-1]] >= a[i] {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i)
		windowMin[i] = a[deque[0]]
	}

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		if l > r {
			l, r = r, l
		}
		curMin := a[l]
		ans := curMin
		for s := l + k; s <= r; s += k {
			segMin := windowMin[s]
			if segMin < curMin {
				curMin = segMin
			}
			ans += curMin
		}
		fmt.Fprintln(writer, ans)
	}
}
