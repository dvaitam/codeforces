package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		target := abs(a[0])
		need := (n - 1) / 2
		count := 0
		for i := 1; i < n; i++ {
			if abs(a[i]) < target {
				count++
			}
		}
		if count >= need {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
