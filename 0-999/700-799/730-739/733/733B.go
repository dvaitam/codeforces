package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	l := make([]int, n)
	r := make([]int, n)
	var sumL, sumR int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &l[i], &r[i])
		sumL += int64(l[i])
		sumR += int64(r[i])
	}
	// current beauty
	diff := sumL - sumR
	if diff < 0 {
		diff = -diff
	}
	best := diff
	bestIdx := 0
	for i := 0; i < n; i++ {
		// swap column i
		newL := sumL - int64(l[i]) + int64(r[i])
		newR := sumR - int64(r[i]) + int64(l[i])
		d := newL - newR
		if d < 0 {
			d = -d
		}
		if d > best {
			best = d
			bestIdx = i + 1
		}
	}
	fmt.Println(bestIdx)
}
