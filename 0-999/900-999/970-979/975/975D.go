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
	var a, b int64
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}

	totalByW := make(map[int64]int64)
	totalByWV := make(map[[2]int64]int64)

	for i := 0; i < n; i++ {
		var x, vx, vy int64
		fmt.Fscan(in, &x, &vx, &vy)
		w := vy - a*vx
		totalByW[w]++
		key := [2]int64{w, vx}
		totalByWV[key]++
	}

	var collisions int64
	for _, cnt := range totalByW {
		collisions += cnt * (cnt - 1) / 2
	}
	for _, cnt := range totalByWV {
		collisions -= cnt * (cnt - 1) / 2
	}

	fmt.Fprintln(out, collisions*2)
}
