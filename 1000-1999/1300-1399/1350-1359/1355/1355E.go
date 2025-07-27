package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var A, R, M int64
	fmt.Fscan(reader, &n, &A, &R, &M)

	heights := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &heights[i])
	}

	sort.Slice(heights, func(i, j int) bool { return heights[i] < heights[j] })

	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + heights[i]
	}

	if M > A+R {
		M = A + R
	}

	low := int64(0)
	high := heights[n-1]

	cost := func(h int64) int64 {
		idx := sort.Search(len(heights), func(i int) bool { return heights[i] > h })
		add := h*int64(idx) - prefix[idx]
		rem := (prefix[n] - prefix[idx]) - h*int64(n-idx)
		if add < 0 {
			add = 0
		}
		if rem < 0 {
			rem = 0
		}
		move := add
		if rem < move {
			move = rem
		}
		return move*M + (add-move)*A + (rem-move)*R
	}

	for high-low > 2 {
		mid1 := low + (high-low)/3
		mid2 := high - (high-low)/3
		if cost(mid1) <= cost(mid2) {
			high = mid2
		} else {
			low = mid1
		}
	}

	ans := cost(low)
	for h := low + 1; h <= high; h++ {
		c := cost(h)
		if c < ans {
			ans = c
		}
	}
	fmt.Fprintln(writer, ans)
}
