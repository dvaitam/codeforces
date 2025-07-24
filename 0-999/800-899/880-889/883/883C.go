package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var f, T, t0 int64
	if _, err := fmt.Fscan(in, &f, &T, &t0); err != nil {
		return
	}
	var a1, t1, p1 int64
	fmt.Fscan(in, &a1, &t1, &p1)
	var a2, t2, p2 int64
	fmt.Fscan(in, &a2, &t2, &p2)

	if f*t0 <= T {
		fmt.Fprintln(out, 0)
		return
	}
	if t1 >= t0 && t2 >= t0 {
		fmt.Fprintln(out, -1)
		return
	}
	// sort by time ascending
	if t1 > t2 {
		a1, a2 = a2, a1
		t1, t2 = t2, t1
		p1, p2 = p2, p1
	}

	const INF int64 = 1<<63 - 1
	ans := INF

	maxN1 := (f + a1 - 1) / a1
	for n1 := int64(0); n1 <= maxN1; n1++ {
		x1 := n1 * a1
		if x1 > f {
			x1 = f
		}
		time1 := x1 * t1
		if time1 > T {
			continue
		}
		remaining := f - x1
		if time1+remaining*t0 <= T {
			cost := n1 * p1
			if cost < ans {
				ans = cost
			}
			continue
		}
		if t2 >= t0 {
			continue
		}
		diff := time1 + remaining*t0 - T
		needBytes := (diff + (t0 - t2) - 1) / (t0 - t2)
		if needBytes > remaining {
			continue
		}
		n2 := (needBytes + a2 - 1) / a2
		cost := n1*p1 + n2*p2
		if cost < ans {
			ans = cost
		}
	}

	if ans == INF {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, ans)
	}
}
