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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	total := int64(n) * int64(n+1) / 2
	for i := 0; i+1 < n; i++ {
		if a[i] != a[i+1] {
			total += int64(i+1) * int64(n-i-1)
		}
	}

	for ; m > 0; m-- {
		var idx int
		var x int64
		fmt.Fscan(in, &idx, &x)
		idx--
		if a[idx] != x {
			if idx-1 >= 0 {
				oldDiff := a[idx-1] != a[idx]
				newDiff := a[idx-1] != x
				if oldDiff != newDiff {
					weight := int64(idx) * int64(n-idx)
					if newDiff {
						total += weight
					} else {
						total -= weight
					}
				}
			}
			if idx+1 < n {
				oldDiff := a[idx] != a[idx+1]
				newDiff := x != a[idx+1]
				if oldDiff != newDiff {
					weight := int64(idx+1) * int64(n-idx-1)
					if newDiff {
						total += weight
					} else {
						total -= weight
					}
				}
			}
			a[idx] = x
		}
		fmt.Fprintln(out, total)
	}
}
