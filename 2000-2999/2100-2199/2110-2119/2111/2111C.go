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

	const inf = int64(1e18)

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		maxRun := make(map[int]int)
		i := 0
		for i < n {
			j := i
			for j < n && arr[j] == arr[i] {
				j++
			}
			runLen := j - i
			if runLen > maxRun[arr[i]] {
				maxRun[arr[i]] = runLen
			}
			i = j
		}

		ans := inf
		for val, run := range maxRun {
			cost := int64(val) * int64(n-run)
			if cost < ans {
				ans = cost
			}
		}
		fmt.Fprintln(out, ans)
	}
}
