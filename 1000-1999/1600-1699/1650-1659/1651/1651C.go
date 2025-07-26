package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
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
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		minA1, minAn := int(1e18), int(1e18)
		minB1, minBn := int(1e18), int(1e18)
		for i := 0; i < n; i++ {
			v := abs(a[0] - b[i])
			if v < minA1 {
				minA1 = v
			}
			v = abs(a[n-1] - b[i])
			if v < minAn {
				minAn = v
			}
		}
		for i := 0; i < n; i++ {
			v := abs(b[0] - a[i])
			if v < minB1 {
				minB1 = v
			}
			v = abs(b[n-1] - a[i])
			if v < minBn {
				minBn = v
			}
		}

		ans := int64(minA1 + minAn + minB1 + minBn)

		candidate := int64(abs(a[0]-b[0]) + abs(a[n-1]-b[n-1]))
		if candidate < ans {
			ans = candidate
		}
		candidate = int64(abs(a[0]-b[n-1]) + abs(a[n-1]-b[0]))
		if candidate < ans {
			ans = candidate
		}
		candidate = int64(abs(a[0]-b[0]) + minAn + minBn)
		if candidate < ans {
			ans = candidate
		}
		candidate = int64(abs(a[n-1]-b[n-1]) + minA1 + minB1)
		if candidate < ans {
			ans = candidate
		}
		candidate = int64(abs(a[0]-b[n-1]) + minAn + minB1)
		if candidate < ans {
			ans = candidate
		}
		candidate = int64(abs(a[n-1]-b[0]) + minA1 + minBn)
		if candidate < ans {
			ans = candidate
		}

		fmt.Fprintln(out, ans)
	}
}
