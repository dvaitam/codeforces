package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n int
		var w int64
		fmt.Fscan(in, &n, &w)

		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		var sum int64
		l, r := n+1, 0 // positions with a[i] < w
		eq := make([]int, 0)
		for i := 0; i < n; i++ {
			sum += a[i]
			if a[i] < w {
				if i+1 < l {
					l = i + 1
				}
				if i+1 > r {
					r = i + 1
				}
			}
			if i+1 < n && a[i] == a[i+1] {
				eq = append(eq, i+1) // pair index (i, i+1) in 1-based
			}
		}

		if len(eq) == 0 {
			fmt.Fprintln(out, sum, 0)
			continue
		}

		if l == n+1 {
			fmt.Fprintln(out, int64(n)*w, 0)
			continue
		}

		target := int64(n)*w - 1
		if sum == target {
			fmt.Fprintln(out, target, 0)
			continue
		}

		if l == r {
			d := int(1e9)
			for _, s := range eq {
				if l > 1 {
					if v := absInt(s - (l - 1)); v < d {
						d = v
					}
				}
				if l <= n-1 {
					if v := absInt(s - l); v < d {
						d = v
					}
				}
			}
			ops := d + 1
			fmt.Fprintln(out, target, ops)
		} else {
			L, R := l, r-1 // interval of pair indices [L, R]
			d := int(1e9)
			for _, s := range eq {
				var dist int
				if s < L {
					dist = L - s
				} else if s > R {
					dist = s - R
				} else {
					dist = 0
				}
				if dist < d {
					d = dist
				}
			}
			ops := (r - l) + d
			fmt.Fprintln(out, target, ops)
		}
	}
}
