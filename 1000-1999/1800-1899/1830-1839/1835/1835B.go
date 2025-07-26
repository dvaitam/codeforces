package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// winCount computes how many target values allow Bytek to win if he picks x.
func winCount(x int64, arr []int64, m int64, k int) int64 {
	n := len(arr)
	// count on the right side of x
	rIdx := sort.Search(n, func(i int) bool { return arr[i] >= x })
	var lenRight int64
	if n-rIdx >= k {
		B := arr[rIdx+k-1]
		rRight := (B - x + 1) / 2
		tRight := x + rRight - 1
		if tRight > m {
			tRight = m
		}
		if tRight >= x {
			lenRight = tRight - x + 1
		}
	} else {
		lenRight = m - x + 1
	}
	// count on the left side of x
	lIdx := sort.Search(n, func(i int) bool { return arr[i] > x }) - 1
	var lenLeft int64
	if lIdx+1 >= k {
		A := arr[lIdx-k+1]
		rLeft := (x - A + 1) / 2
		tLeft := x - rLeft + 1
		if tLeft < 0 {
			tLeft = 0
		}
		if tLeft <= x {
			lenLeft = x - tLeft + 1
		}
	} else {
		lenLeft = x + 1
	}
	total := lenLeft + lenRight
	if lenLeft > 0 && lenRight > 0 {
		total--
	}
	return total
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	var m int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

	candMap := make(map[int64]struct{})
	add := func(v int64) {
		if v >= 0 && v <= m {
			candMap[v] = struct{}{}
		}
	}
	add(0)
	add(m)
	for i := 0; i < n; i++ {
		add(arr[i])
		add(arr[i] - 1)
		add(arr[i] + 1)
		if i < n-1 {
			mid := (arr[i] + arr[i+1]) / 2
			add(mid)
			add(mid - 1)
			add(mid + 1)
		}
	}
	candidates := make([]int64, 0, len(candMap))
	for v := range candMap {
		candidates = append(candidates, v)
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i] < candidates[j] })

	var bestX int64
	var bestVal int64 = -1
	for _, x := range candidates {
		val := winCount(x, arr, m, k)
		if val > bestVal || (val == bestVal && x < bestX) {
			bestVal = val
			bestX = x
		}
	}
	fmt.Fprintf(writer, "%d %d\n", bestVal, bestX)
}
