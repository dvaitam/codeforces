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
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	left := make(map[int64]int64)
	right := make(map[int64]int64)
	for _, v := range a {
		right[v]++
	}

	var ans int64
	for _, v := range a {
		right[v]--
		if v%k == 0 {
			ans += left[v/k] * right[v*k]
		}
		left[v]++
	}

	fmt.Fprintln(out, ans)
}
