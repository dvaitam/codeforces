package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var a, b, k int
	if _, err := fmt.Fscan(in, &n, &a, &b, &k); err != nil {
		return
	}
	h := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &h[i])
	}

	cycle := a + b
	need := make([]int, 0, n)
	for _, hp := range h {
		r := hp % cycle
		if r == 0 {
			r = cycle
		}
		extra := (r - 1) / a
		if extra > 0 {
			need = append(need, extra)
		} else {
			need = append(need, 0)
		}
	}

	sort.Ints(need)

	points := 0
	for _, v := range need {
		if k < v {
			break
		}
		k -= v
		points++
	}

	fmt.Fprintln(out, points)
}
