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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	if n <= 2 {
		fmt.Fprintln(out, 0)
		return
	}

	const INF = int(1e9)
	res := INF

	deltas := []int{-1, 0, 1}
	for _, d1 := range deltas {
		for _, d2 := range deltas {
			start := arr[0] + d1
			diff := (arr[1] + d2) - start
			changes := 0
			if d1 != 0 {
				changes++
			}
			if d2 != 0 {
				changes++
			}
			ok := true
			for i := 2; i < n; i++ {
				expected := start + diff*i
				delta := expected - arr[i]
				if delta < -1 || delta > 1 {
					ok = false
					break
				}
				if delta != 0 {
					changes++
				}
			}
			if ok && changes < res {
				res = changes
			}
		}
	}

	if res == INF {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, res)
	}
}
