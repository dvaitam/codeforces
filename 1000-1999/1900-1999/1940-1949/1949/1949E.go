package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	h := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &h[i])
	}

	bestX, bestY := 1, k-1
	bestTime := math.Inf(1)
	for x := 1; x <= k-1; x++ {
		y := k - x
		var hits int64
		for _, v := range h {
			hits += (v + int64(x) - 1) / int64(x)
		}
		t := float64(hits) / float64(y)
		if t < bestTime {
			bestTime = t
			bestX, bestY = x, y
		}
	}

	fmt.Fprintf(out, "%d %d\n", bestX, bestY)
}
