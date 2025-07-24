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

	var n int
	var t int64
	if _, err := fmt.Fscan(in, &n, &t); err != nil {
		return
	}
	g := make([]int64, n)
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &g[i], &c[i])
	}
	d := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &d[i])
	}

	ans := int64(math.MaxInt64)
	for start := int64(0); start < t; start++ {
		time := start
		for i := 0; i < n; i++ {
			offset := (c[i] + time) % t
			if offset >= g[i] {
				time += t - offset
			}
			if i < n-1 {
				time += d[i]
			}
		}
		driveTime := time - start
		if driveTime < ans {
			ans = driveTime
		}
	}
	fmt.Fprintln(out, ans)
}
