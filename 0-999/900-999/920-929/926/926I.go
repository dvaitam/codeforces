package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	times := make([]int, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		var h, m int
		fmt.Sscanf(s, "%d:%d", &h, &m)
		times[i] = h*60 + m
	}
	sort.Ints(times)
	maxGap := 0
	for i := 0; i < n; i++ {
		next := (i + 1) % n
		endPrev := times[i] + 1
		startNext := times[next]
		if next == 0 {
			startNext += 1440
		}
		gap := startNext - endPrev
		if gap > maxGap {
			maxGap = gap
		}
	}
	fmt.Printf("%02d:%02d\n", maxGap/60, maxGap%60)
}
