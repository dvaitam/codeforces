package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		d := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &d[i])
		}
		l := make([]int, n+1)
		r := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &l[i], &r[i])
		}

		mn := make([]int, n+1)
		mx := make([]int, n+1)
		mn[0], mx[0] = 0, 0
		ok := true
		for i := 1; i <= n && ok; i++ {
			low, high := 0, 0
			switch d[i] {
			case 0:
				low, high = mn[i-1], mx[i-1]
			case 1:
				low, high = mn[i-1]+1, mx[i-1]+1
			default: // -1
				low, high = mn[i-1], mx[i-1]+1
			}
			low = max(low, l[i])
			high = min(high, r[i])
			if low > high {
				ok = false
				break
			}
			mn[i], mx[i] = low, high
		}
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}

		gl := make([]int, n+1)
		gr := make([]int, n+1)
		gl[n], gr[n] = mn[n], mx[n]
		for i := n; i >= 1; i-- {
			low, high := 0, 0
			switch d[i] {
			case 0:
				low, high = gl[i], gr[i]
			case 1:
				low, high = gl[i]-1, gr[i]-1
			default:
				low, high = gl[i]-1, gr[i]
			}
			if low < 0 {
				low = 0
			}
			low = max(low, mn[i-1])
			high = min(high, mx[i-1])
			if low > high {
				ok = false
				break
			}
			gl[i-1], gr[i-1] = low, high
		}
		if !ok || !(gl[0] <= 0 && 0 <= gr[0]) {
			fmt.Fprintln(out, -1)
			continue
		}

		h := make([]int, n+1)
		h[0] = 0
		res := make([]int, n+1)
		for i := 1; i <= n; i++ {
			prev := h[i-1]
			switch d[i] {
			case 0:
				h[i] = prev
				if h[i] < gl[i] || h[i] > gr[i] {
					ok = false
				}
			case 1:
				h[i] = prev + 1
				if h[i] < gl[i] || h[i] > gr[i] {
					ok = false
				}
			default:
				if prev >= gl[i] && prev <= gr[i] {
					h[i] = prev
				} else if prev+1 >= gl[i] && prev+1 <= gr[i] {
					h[i] = prev + 1
				} else {
					ok = false
				}
			}
			res[i] = h[i] - prev
		}
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, res[i])
		}
		fmt.Fprintln(out)
	}
}
