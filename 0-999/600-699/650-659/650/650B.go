package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func upperBound(arr []int64, val int64) int {
	// returns first index i such that arr[i] > val
	return sort.Search(len(arr), func(i int) bool { return arr[i] > val })
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, a, b int
	var T int64
	if _, err := fmt.Fscan(in, &n, &a, &b, &T); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	cost := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if s[i-1] == 'w' {
			cost[i] = int64(1 + b)
		} else {
			cost[i] = 1
		}
	}
	first := cost[1]

	right := make([]int64, n)
	left := make([]int64, n)
	right[0] = first
	left[0] = first
	for i := 1; i < n; i++ {
		right[i] = right[i-1] + int64(a) + cost[i+1]
		left[i] = left[i-1] + int64(a) + cost[n-i+1]
	}

	maxPhoto := 0
	for i := 0; i < n; i++ {
		if right[i] <= T && i+1 > maxPhoto {
			maxPhoto = i + 1
		}
		if left[i] <= T && i+1 > maxPhoto {
			maxPhoto = i + 1
		}
	}

	for r := 1; r < n; r++ {
		if right[r] > T {
			break
		}
		remain := T - right[r] - int64(r*a)
		if remain < 0 {
			continue
		}
		limit := remain + first
		l := upperBound(left, limit) - 1
		if l > n-1-r {
			l = n - 1 - r
		}
		if r+1+l > maxPhoto {
			maxPhoto = r + 1 + l
		}
	}

	for l := 1; l < n; l++ {
		if left[l] > T {
			break
		}
		remain := T - left[l] - int64(l*a)
		if remain < 0 {
			continue
		}
		limit := remain + first
		r := upperBound(right, limit) - 1
		if r > n-1-l {
			r = n - 1 - l
		}
		if l+1+r > maxPhoto {
			maxPhoto = l + 1 + r
		}
	}

	fmt.Println(maxPhoto)
}
